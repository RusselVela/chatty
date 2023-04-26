package service

import (
	"context"
	"github.com/RusselVela/chatty/internal/app/domain"
	"github.com/gorilla/websocket"
	"github.com/knadh/koanf"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

// WebsocketHandler is used for managing websocket connections
type WebsocketHandler struct {
	wsUpgrader       websocket.Upgrader
	wsWriteWait      time.Duration
	wsPongWait       time.Duration
	wsPingPeriod     time.Duration
	wsMaxMessageSize int64
}

var (
	errConnectionAlreadyExists = &domain.ErrorCode{Code: 103, Message: "a connection already exists"}
	errUpgradingWebsocket      = &domain.ErrorCode{Code: 109, Message: "failed to upgrade to a websocket connection: %s"}
)
var wsHandler *WebsocketHandler

// NewWebsocketHandler creates a new WebsocketHandler
func NewWebsocketHandler(k *koanf.Koanf) *WebsocketHandler {
	if wsHandler != nil {
		return wsHandler
	}

	wsHandler = &WebsocketHandler{
		wsUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}

	// Time allowed to write a message to the peer.
	wsHandler.wsWriteWait = time.Duration(k.Int("ws.writewait")) * time.Second

	// Time allowed to read the next pong message from the peer.
	wsHandler.wsPongWait = time.Duration(k.Int("ws.pongwait")) * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	wsHandler.wsPingPeriod = (wsHandler.wsPongWait * 9) / 10

	// Maximum message size allowed from peer.
	wsHandler.wsMaxMessageSize = k.Int64("ws.maxmessagesize")

	return wsHandler
}

// UpgradeConnection will upgrade an existing connection for the supplied websocket client.
func (wh *WebsocketHandler) UpgradeConnection(ctx echo.Context, wsClient *WsClient) error {
	if wsClient.wsConn != nil {
		return ctx.JSON(http.StatusConflict, errConnectionAlreadyExists)
	}

	wsConn, err := wsHandler.wsUpgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return ctx.JSON(http.StatusPreconditionFailed, errUpgradingWebsocket.Clone(err.Error()))
	}
	wsClient.ctx, wsClient.cancel = context.WithCancel(context.Background())
	wsClient.wsConn = wsConn
	wsClient.wsHandler = wh

	return nil
}
