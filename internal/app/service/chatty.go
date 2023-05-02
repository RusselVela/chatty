package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	chattyredis "github.com/RusselVela/chatty/internal/app/datasourcce/repository/redis"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
)

var (
	errUserExists            = &ErrorCode{Status: http.StatusBadRequest, Code: 201, Message: "username %s already taken"}
	errUserNotExist          = &ErrorCode{Status: http.StatusBadRequest, Code: 202, Message: "username %s does not exist"}
	errInvalidCredentials    = &ErrorCode{Status: http.StatusBadRequest, Code: 203, Message: "invalid credentials"}
	errTokenGeneration       = &ErrorCode{Status: http.StatusInternalServerError, Code: 204, Message: "error generating JWT token: %s"}
	errChannelExists         = &ErrorCode{Status: http.StatusBadRequest, Code: 301, Message: "channel %s already exists"}
	errChannelNotExist       = &ErrorCode{Status: http.StatusBadRequest, Code: 302, Message: "channel %s does not exist"}
	errChannelCreation       = &ErrorCode{Status: http.StatusBadRequest, Code: 303, Message: "channel creation failed: %s"}
	errUserAlreadySubscribed = &ErrorCode{Status: http.StatusBadRequest, Code: 305, Message: "user %s already is a member of channel %s"}
	errUserClientNotExist    = &ErrorCode{Status: http.StatusInternalServerError, Code: 401, Message: "the client associated to user %s does not exist"}
)

// ChattyService implements the ChattyService interface to be consumed by the web layer
type ChattyService struct {
	redisClient *chattyredis.Client
}

// NewChattyService returns a struct that will handle all service functionality
func NewChattyService(redisClient *chattyredis.Client) *ChattyService {
	return &ChattyService{
		redisClient: redisClient,
	}
}

// Signup creates new users on the service
func (cs *ChattyService) Signup(username string, password string) (string, string, error) {
	user := inmemory.GetUserByName(username)
	if user != nil {
		return "", "", errUserExists.Clone(username)
	}

	user, err := inmemory.NewUser(username, password)
	if err != nil {
		return "", "", errUserExists.Clone(username)
	}

	zap.S().Infof("new user %s created: %s", user.Id.String(), user.Username)
	return user.Id.String(), user.Username, nil
}

// Login authenticates an existing user. Returns a token that can be used for further operations
func (cs *ChattyService) Login(username string, password string) (string, error) {
	user := inmemory.GetUserByName(username)
	if user == nil || (user.Username != username || user.Password != password) {
		return "", errInvalidCredentials
	}
	zap.S().Infof("user login: %s", username)

	token, err := generateJWT(*user, 0)
	if err != nil {
		return "", errTokenGeneration.Clone(err.Error())
	}

	return token, nil
}

// HandleConnection takes a http request and upgrades it to a websocket connection. The socket created here will be used
// to communicate with the client
func (cs *ChattyService) HandleConnection(ctx echo.Context) error {
	wsClient := &UserClient{}
	err := wsHandler.UpgradeConnection(ctx, wsClient)
	if err != nil {
		return err
	}

	userId, err := wsClient.readAuthMessage()
	if err != nil {
		wsClient.message <- domain.Message{
			Text: "invalid authentication token",
		}

		if err = wsClient.wsConn.Close(); err != nil {
			zap.S().Warnf("ws upgrade: ws connection close error %v", err)
		}
		wsClient.cancel()
		wsClient.release()
		zap.S().Infof("Returning WS auth error")
		return err
	}

	user := inmemory.GetUser(userId)
	if user == nil {
		msg := "user %s not found"
		zap.L().Error(fmt.Sprintf(msg, user.Id.String()))
		return errUserNotExist.Clone(user.Id.String())
	}

	user.Online = true
	wsClient.user = user
	clients[user.Id.String()] = wsClient

	err = cs.sendPreviousMessages(wsClient)
	if err != nil {
		zap.S().Errorf("failed to send previous messages: %s", err.Error())
	}
	go wsClient.writeMessage()
	wsClient.readMessages()

	wsClient.ctx, wsClient.cancel = nil, nil
	wsClient.wsConn = nil

	return nil
}

// GetConnectionToken returns a short-lived token that must be sent as first message after connection
// with the server has been established
func (cs *ChattyService) GetConnectionToken(ctx echo.Context, userId string) (string, error) {
	user := inmemory.GetUser(userId)
	if user == nil {
		return "", errInvalidCredentials
	}
	zap.S().Infof("generating WS auth token for user: %s", user.Id)

	token, err := generateJWT(*user, 10)
	if err != nil {
		return "", errTokenGeneration.Clone(err.Error())
	}

	inmemory.AddTokenToUser(user.Id.String(), token)

	return token, nil
}

// CreateChannel creates a new group for multiple users to interact
func (cs *ChattyService) CreateChannel(name string, visibility string, ownerId string) (*inmemory.ChannelBean, error) {
	user := inmemory.GetUser(ownerId)
	if user == nil {
		return nil, errUserNotExist.Clone(ownerId)
	}

	channel := inmemory.GetChannelByName(name)
	if channel != nil {
		return nil, errChannelExists.Clone(name)
	}

	channel, err := inmemory.NewChannel(name, ownerId, visibility)
	if err != nil {
		return nil, errChannelCreation.Clone(err)
	}

	channelClient := NewChannelClient(channel)
	go channelClient.Start()

	zap.S().Infof("User %s created Channel: %s", user.Username, channel.Name)
	return channel, nil
}

// SubscribeChannel subscribes the given userId to channelId. After this action, all messages sent to channelId
// will also be sent to userId
func (cs *ChattyService) SubscribeChannel(userId string, channelId string) error {
	channel := inmemory.GetChannel(channelId)
	if channel == nil {
		return errChannelNotExist.Clone(channelId)
	}

	user := inmemory.GetUser(userId)
	if user == nil {
		return errUserNotExist.Clone(userId)
	}

	if _, found := channel.Members[user.Id.String()]; found {
		zap.S().Infof("User %s already a member of channel %s", user.Id, channel.Id)
		return errUserAlreadySubscribed.Clone(user.Username, channel.Name)
	}

	channelClient, found := channelClients[channel.Id.String()]
	if !found {
		// Channel client not placed for some reason. Start it
		channelClient = NewChannelClient(channel)
		go channelClient.Start()
		zap.S().Infof("Client for Channel %s started", channel.Name)
	}

	channelClient.Subscribe <- user
	user.Subscriptions = append(user.Subscriptions, channel.Id.String())
	zap.S().Infof("User %s joined Channel: %s", user.Username, channel.Name)

	return nil
}

// GetUsers returns a list of all users in the service, online or not.
func (cs *ChattyService) GetUsers() ([]*inmemory.UserBean, error) {
	return inmemory.GetUsers(), nil
}

// GetChannels returns a list of all public channels.
func (cs *ChattyService) GetChannels() ([]*inmemory.ChannelBean, error) {
	return inmemory.GetChannels(), nil
}

// sendPreviousMessages sends all messages for a user that were delivered while offline
func (cs *ChattyService) sendPreviousMessages(wsClient *UserClient) error {
	messages := cs.redisClient.GetMessages(wsClient.user.Id.String())
	for _, msg := range messages {
		wsClient.message <- *msg
	}
	_, err := cs.redisClient.ClearMessageQueue(wsClient.user.Id.String())
	if err != nil {
		return err
	}

	return nil
}
