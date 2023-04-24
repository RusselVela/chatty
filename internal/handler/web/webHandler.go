package web

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=oapi-cfg.yaml ../../../api/chatty-service-api.yaml

type loggerContextKeyType struct{}

var loggerContextKey = loggerContextKeyType{}

type ChattyService interface {
	Login(username string, password string) (string, error)
	PostMessage(token string, recipient string, message string) (int64, error)
}

type WebHandler struct {
	service ChattyService
}

func NewWebHandler(service ChattyService) *WebHandler {
	return &WebHandler{
		service: service,
	}
}

func (wh *WebHandler) PublicPostLogin(ctx echo.Context) error {
	l := L(ctx.Request().Context())
	l.Info("Posting a new Message")
	return nil
}

func (wh *WebHandler) PublicPostMessage(ctx echo.Context) error {
	l := L(ctx.Request().Context())
	l.Info("Posting a new Message")
	return nil
}

// L retrieves logger value from context
func L(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerContextKey).(*zap.Logger); ok {
		return logger
	}
	return nil
}
