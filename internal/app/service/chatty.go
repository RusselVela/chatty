package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
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
}

func NewChattyService() *ChattyService {
	return &ChattyService{}
}

func (cs *ChattyService) Signup(username string, password string) (string, string, error) {
	user := inmemory.Users.Get(username)
	if user != nil {
		return "", "", errUserExists.Clone(username)
	}

	user, err := inmemory.Users.NewUser(username, password)
	if err != nil {
		return "", "", errUserExists.Clone(username)
	}

	zap.L().Info("new user created")
	return user.Id, user.Username, nil
}

func (cs *ChattyService) Login(username string, password string) (string, error) {
	user := inmemory.Users.Get(username)
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

func (cs *ChattyService) HandleConnections(ctx echo.Context, username string) error {
	user := inmemory.Users.Get(username)
	if user == nil {
		msg := "user %s not found"
		zap.L().Error(fmt.Sprintf(msg, username))
		return errUserNotExist.Clone(username)
	}

	wsClient := &UserClient{}
	err := wsHandler.UpgradeConnection(ctx, wsClient)
	if err != nil {
		return err
	}
	wsClient.user = user

	clients[user.Username] = wsClient

	wsClient.readMessages()

	wsClient.ctx, wsClient.cancel = nil, nil
	wsClient.wsConn = nil

	return nil
}

func (cs *ChattyService) CreateChannel(name string, visibility string, owner string) error {
	user := inmemory.Users.Get(owner)
	if user == nil {
		return errUserNotExist.Clone(owner)
	}

	channel := inmemory.Channels.Get(name)
	if channel != nil {
		return errChannelExists.Clone(name)
	}

	channel, err := inmemory.Channels.NewChannel(name, owner, visibility)
	if err != nil {
		return errChannelCreation.Clone(err)
	}

	channelClient := NewChannelClient(channel)
	go channelClient.Start()

	zap.S().Infof("User %s created Channel: %s", user.Username, channel.Name)
	return nil
}

func (cs *ChattyService) SubscribeChannel(username string, channelName string) error {
	channel := inmemory.Channels.Get(channelName)
	if channel == nil {
		return errChannelNotExist.Clone(channelName)
	}

	user := inmemory.Users.Get(username)
	if user == nil {
		return errUserNotExist.Clone(username)
	}

	if _, found := channel.Members[user.Username]; found {
		zap.S().Infof("User %s already a member of channel %s", user.Username, channel.Name)
		return errUserAlreadySubscribed.Clone(user.Username, channel.Name)
	}

	channelClient, found := channelClients[channel.Name]
	if !found {
		// Channel client not placed for some reason. Start it
		channelClient = NewChannelClient(channel)
		go channelClient.Start()
		zap.S().Infof("Client for Channel %s started", channel.Name)
	}

	channelClient.Subscribe <- user
	user.Subscriptions = append(user.Subscriptions, channel.Name)
	zap.S().Infof("User %s joined Channel: %s", user.Username, channel.Name)

	return nil
}
