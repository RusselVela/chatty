package service

import (
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ChattyService struct {
}

func NewChattyService() *ChattyService {
	return &ChattyService{}
}

func (cs *ChattyService) Signup(username string, password string) (string, string, error) {
	user := inmemory.Users.Get(username)
	if user != nil {
		return "", "", fmt.Errorf("user %s already exist", username)
	}

	err, user := inmemory.Users.NewUser(username, password)
	if err != nil {
		return "", "", fmt.Errorf("error creating new user: %w", err)
	}

	zap.L().Info("new user created")
	return user.Id, user.Username, nil
}

func (cs *ChattyService) Login(username string, password string) (string, error) {
	user := inmemory.Users.Get(username)
	if user == nil || (user.Username != username || user.Password != password) {
		return "", fmt.Errorf("invalid credentials")
	}
	zap.L().Info(fmt.Sprintf("user login: %s", username))

	token, err := generateJWT(*user)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (cs *ChattyService) HandleConnections(ctx echo.Context, username string) error {
	user := inmemory.Users.Get(username)
	if user == nil {
		msg := "user %s not found"
		zap.L().Error(fmt.Sprintf(msg, username))
		return fmt.Errorf(msg, username)
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
