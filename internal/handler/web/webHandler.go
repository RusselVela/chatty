package web

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=oapi-cfg.yaml ../../../api/chatty-service-api.yaml

type loggerContextKeyType struct{}

var loggerContextKey = loggerContextKeyType{}

type ChattyService interface {
	Signup(username string, password string) (string, string, error)
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

func (wh *WebHandler) PublicPostSignup(ctx echo.Context) error {
	request := PublicPostSignupJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, N401FailedLogin{
			Ok:    false,
			Error: err.Error(),
		})
	}

	id, username, err := wh.service.Signup(request.Username, request.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, N401FailedLogin{
			Ok:    false,
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, N201SuccessfulSignUp{
		Id:       id,
		Username: username,
		Ok:       true,
	})
}

func (wh *WebHandler) PublicPostLogin(ctx echo.Context) error {
	request := PublicPostLoginJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, N401FailedLogin{
			Ok:    false,
			Error: err.Error(),
		})
	}

	token, err := wh.service.Login(request.Username, request.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, N401FailedLogin{
			Ok:    false,
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, N200SuccessfulLogin{
		Ok:    true,
		Token: token,
	})
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
