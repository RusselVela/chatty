package service

import (
	"context"
	"fmt"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type UserClient struct {
	user      *inmemory.UserBean
	wsConn    *websocket.Conn
	wsHandler *WebsocketHandler
	ctx       context.Context
	cancel    context.CancelFunc
}

// readAuthMessage waits for the first message that comes from client. Then parses it to retrieve a token.
// If no token is sent, it returns an error
func (uc *UserClient) readAuthMessage() (string, error) {
	var msg domain.Message

	for {
		err := uc.wsConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Warn("ws read error: %v", err)
			}
			return "", fmt.Errorf("invalid authentication handshake: %w", err)
		}
		zap.S().Infof("received auth message: %s", msg.Text)
		break
	}

	tokenStr := msg.Text
	token, err := parseJWT(tokenStr)
	if err != nil {
		return "", fmt.Errorf("invalid authentication handshake: %w", err)
	}

	userId := token.Claims.(*JWTCustomClaims).Id

	cachedToken := inmemory.GetToken(userId)
	if cachedToken == "" || cachedToken != tokenStr {
		return "", fmt.Errorf("invalid auth token")
	}

	// token has been used, remove it
	inmemory.DeleteTokenToUser(userId)

	return userId, nil
}

// readMessages receives all incoming messages from client and sends them to the broadcaster to be properly delivered to
// their destination
func (uc *UserClient) readMessages() {
	defer func() {
		if err := uc.wsConn.Close(); err != nil {
			zap.S().Warnf("ws upgrade: ws connection close error %v", err)
		}
		uc.cancel()
		uc.release()
	}()

	for {
		var msg domain.Message
		err := uc.wsConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Warn("ws read error: %+v", err)
			}
			return
		}
		zap.S().Infof("received message: %v", msg)

		msg.SourceId = uc.user.Id.String()
		broadcaster <- msg
	}
}

// writeMessage sends a domain.Message object to the client of this websocket
func (uc *UserClient) writeMessage(msg domain.Message) {
	err := uc.wsConn.WriteJSON(msg)
	if err != nil && !websocket.IsCloseError(err, websocket.CloseGoingAway) {
		removeClient(uc.user.Id.String())
	}
}

// release clears all the resources no longer needed by this websocket client. Normally invoked after closing the socket
func (uc *UserClient) release() {
	uc.ctx = nil
	uc.cancel = nil
	uc.wsConn = nil
	uc.wsHandler = nil
	if uc.user != nil {
		removeClient(uc.user.Id.String())
		uc.user.Online = false
		uc.user = nil
	}
}
