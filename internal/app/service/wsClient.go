package service

import (
	"context"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WsClient struct {
	user      *inmemory.UserBean
	wsConn    *websocket.Conn
	wsHandler *WebsocketHandler
	ctx       context.Context
	cancel    context.CancelFunc
}

func (wsc *WsClient) readMessages() {
	defer func() {
		if err := wsc.wsConn.Close(); err != nil {
			zap.S().Warnf("ws upgrade: ws connection close error %v", err)
		}
		wsc.cancel()
		wsc.release()
	}()

	for {
		var msg domain.Message
		err := wsc.wsConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Warn("ws read error: %+v", err)
			}
			return
		}
		zap.S().Infof("received message: %v", msg)

		msg.Sender = wsc.user.Username
		broadcaster <- msg
	}
}

func (wsc *WsClient) writeMessage(msg domain.Message) {
	err := wsc.wsConn.WriteJSON(msg)
	if err != nil && !websocket.IsCloseError(err, websocket.CloseGoingAway) {
		removeClient(wsc.user.Username)
	}
}

func (wsc *WsClient) release() {
	wsc.ctx = nil
	wsc.cancel = nil
	wsc.wsConn = nil
	removeClient(wsc.user.Username)
}
