package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	chattyredis "github.com/RusselVela/chatty/internal/app/datasourcce/repository/redis"
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

type ChattyService struct {
	redisClient *chattyredis.Client
}

func NewChattyService(redisClient *chattyredis.Client) *ChattyService {
	return &ChattyService{
		redisClient: redisClient,
	}
}

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

func (cs *ChattyService) Login(username string, password string) (string, error) {
	user := inmemory.GetUserByName(username)
	if user == nil || (user.Username != username || user.Password != password) {
		return "", errInvalidCredentials
	}
	zap.S().Infof("user login: %s", username)

	token, err := generateJWT(*user)
	if err != nil {
		return "", errTokenGeneration.Clone(err.Error())
	}

	return token, nil
}

func (cs *ChattyService) HandleConnections(ctx echo.Context, userId string) error {
	user := inmemory.GetUser(userId)
	if user == nil {
		msg := "user %s not found"
		zap.L().Error(fmt.Sprintf(msg, user.Id.String()))
		return errUserNotExist.Clone(user.Id.String())
	}

	wsClient := &UserClient{}
	err := wsHandler.UpgradeConnection(ctx, wsClient)
	if err != nil {
		return err
	}
	wsClient.user = user

	clients[user.Id.String()] = wsClient

	err = cs.sedPreviousMessages(wsClient)
	if err != nil {
		zap.S().Errorf("failed to send previous messages: %s", err.Error())
	}
	wsClient.readMessages()

	wsClient.ctx, wsClient.cancel = nil, nil
	wsClient.wsConn = nil

	return nil
}

func (cs *ChattyService) CreateChannel(name string, visibility string, ownerId string) error {
	user := inmemory.GetUser(ownerId)
	if user == nil {
		return errUserNotExist.Clone(ownerId)
	}

	channel := inmemory.GetChannelByName(name)
	if channel != nil {
		return errChannelExists.Clone(name)
	}

	channel, err := inmemory.NewChannel(name, ownerId, visibility)
	if err != nil {
		return errChannelCreation.Clone(err)
	}

	channelClient := NewChannelClient(channel)
	go channelClient.Start()

	zap.S().Infof("User %s created Channel: %s", user.Username, channel.Name)
	return nil
}

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

func (cs *ChattyService) sedPreviousMessages(wsClient *UserClient) error {
	messages := cs.redisClient.GetMessages(wsClient.user.Id.String())
	for _, msg := range messages {
		wsClient.writeMessage(*msg)
	}
	_, err := cs.redisClient.ClearMessageQueue(wsClient.user.Id.String())
	if err != nil {
		return err
	}

	return nil
}
