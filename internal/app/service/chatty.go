package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

var (
	errUserExists         = &domain.ErrorCode{Code: 201, Message: "a User already exists with the given name"}
	errUserNotExist       = &domain.ErrorCode{Code: 202, Message: "the specified User does not exist"}
	errInvalidCredentials = &domain.ErrorCode{Code: 203, Message: "invalid credentials"}
	errTokenGeneration    = &domain.ErrorCode{Code: 204, Message: "error generating JWT token"}

	errChannelExists         = &domain.ErrorCode{Code: 301, Message: "a Channel already exists with the given name"}
	errChannelNotExist       = &domain.ErrorCode{Code: 302, Message: "the specified Channel does not exist"}
	errChannelCreation       = &domain.ErrorCode{Code: 303, Message: "channel creation failed"}
	errUserAlreadySubscribed = &domain.ErrorCode{Code: 305, Message: "the User already is a member of the Channel"}

	errUserClientNotExist = &domain.ErrorCode{Code: 401, Message: "the wsClient associated to the user does not exist"}
)

type ChattyService struct {
}

func NewChattyService() *ChattyService {
	return &ChattyService{}
}

func (cs *ChattyService) Signup(username string, password string) (string, string, error) {
	user := inmemory.Users.Get(username)
	if user != nil {
		return "", "", errUserExists.Clone()
	}

	user, err := inmemory.Users.NewUser(username, password)
	if err != nil {
		return "", "", errUserExists.Clone()
	}

	zap.L().Info("new user created")
	return user.Id, user.Username, nil
}

func (cs *ChattyService) Login(username string, password string) (string, error) {
	user := inmemory.Users.Get(username)
	if user == nil || (user.Username != username || user.Password != password) {
		return "", errInvalidCredentials.Clone()
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
		return errUserNotExist
	}

	wsClient := &WsClient{}
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
		return errUserNotExist
	}

	channel := inmemory.Channels.Get(name)
	if channel != nil {
		return errChannelExists.Clone()
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
		return errChannelNotExist.Clone()
	}

	user := inmemory.Users.Get(username)
	if user == nil {
		return errUserNotExist
	}

	if _, found := channel.Members[user.Username]; found {
		zap.S().Infof("User %s already a member of channel %s", user.Username, channel.Name)
		return errUserAlreadySubscribed
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
