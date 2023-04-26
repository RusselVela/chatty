package web

import (
	"github.com/RusselVela/chatty/internal/app/service"
	"github.com/golang-jwt/jwt/v4"
	"net/http"

	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=oapi-cfg.yaml ../../../api/chatty-service-api.yaml

type loggerContextKeyType struct{}

var loggerContextKey = loggerContextKeyType{}

type ChattyService interface {
	Signup(username string, password string) (string, string, error)
	Login(username string, password string) (string, error)
	CreateChannel(name string, visibility string, owner string) error
	SubscribeChannel(username string, channelName string) error
	HandleConnections(ctx echo.Context, token string) error
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

func (wh *WebHandler) PublicPostToken(ctx echo.Context) error {
	request := PublicPostTokenJSONRequestBody{}
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

func (wh *WebHandler) PublicGetWs(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	err := wh.service.HandleConnections(ctx, claims.Username)
	if err != nil {

	}
	return nil
}

func (wh *WebHandler) PublicPostChannels(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	request := PublicPostChannelsJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, N400FailedChannelCreation{
			Ok:    false,
			Error: err.Error(),
		})
	}

	err := wh.service.CreateChannel(request.Name, request.Type, claims.Username)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, N500InternalServerError{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, N201SuccessChannelCreation{
		Ok:   true,
		Name: request.Name,
	})
}

func (wh *WebHandler) PublicPostChannelsSubscribe(ctx echo.Context, name string) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	err := wh.service.SubscribeChannel(claims.Username, name)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, N500InternalServerError{
			Error: err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, N200SuccessChannelSubscribe{
		Ok:   true,
		Name: name,
	})
}
