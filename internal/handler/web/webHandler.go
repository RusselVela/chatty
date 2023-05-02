package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"
	"github.com/RusselVela/chatty/internal/app/service"
	"github.com/golang-jwt/jwt/v4"

	"github.com/labstack/echo/v4"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=oapi-cfg.yaml ../../../api/chatty-service-api.yaml

type ChattyService interface {
	Signup(username string, password string) (string, string, error)
	Login(username string, password string) (string, error)
	GetUsers() ([]*inmemory.UserBean, error)
	CreateChannel(name string, visibility string, owner string) (*inmemory.ChannelBean, error)
	SubscribeChannel(username string, channelName string) error
	GetChannels() ([]*inmemory.ChannelBean, error)
	HandleConnection(ctx echo.Context) error
	GetConnectionToken(ctx echo.Context, userId string) (string, error)
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
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	id, username, err := wh.service.Signup(request.Username, request.Password)
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	return ctx.JSON(http.StatusCreated, N201SuccessSignUp{
		Id:       id,
		Username: username,
	})
}

func (wh *WebHandler) PublicPostToken(ctx echo.Context) error {
	request := PublicPostTokenJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	token, err := wh.service.Login(request.Username, request.Password)
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	return ctx.JSON(http.StatusOK, N200SuccessLogin{
		Token: token,
	})
}

func (wh *WebHandler) PublicGetWs(ctx echo.Context) error {
	err := wh.service.HandleConnection(ctx)
	if err != nil {
		return nil
	}
	return nil
}

func (wh *WebHandler) PublicGetWsToken(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	token, err := wh.service.GetConnectionToken(ctx, claims.Id)
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	return ctx.JSON(http.StatusOK, N200SuccessLogin{
		Token: token,
	})
}

func (wh *WebHandler) PublicPostChannels(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	request := PublicPostChannelsJSONRequestBody{}
	if err := ctx.Bind(&request); err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	channel, err := wh.service.CreateChannel(request.Name, request.Type, claims.Id)
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	return ctx.JSON(http.StatusCreated, N201SuccessChannelCreation{
		Id:      channel.Id.String(),
		OwnerId: channel.OwnerId,
		Name:    request.Name,
	})
}

func (wh *WebHandler) PublicPostChannelsSubscribe(ctx echo.Context, id string) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	err := wh.service.SubscribeChannel(claims.Id, id)
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}

	return ctx.JSON(http.StatusOK, N200SuccessChannelSubscribe{
		Ok:   true,
		Name: id,
	})
}

func (wh *WebHandler) PublicGetChannels(ctx echo.Context) error {
	channelBeans, err := wh.service.GetChannels()
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}
	channels := make([]Channel, 0)
	for _, cb := range channelBeans {
		channel := wh.beanToChannel(cb)
		channels = append(channels, channel)
	}
	list := SuccessGetChannels{Channels: channels}
	return ctx.JSON(http.StatusOK, list)
}

func (wh *WebHandler) PublicGetUsers(ctx echo.Context) error {
	userBeans, err := wh.service.GetUsers()
	if err != nil {
		status, errMsg := wh.toErrorMessage(err)
		return ctx.JSON(status, errMsg)
	}
	users := make([]User, 0)
	for _, ub := range userBeans {
		user := wh.beanToUser(ub)
		users = append(users, user)
	}
	list := SuccessGetUsers{Users: users}
	return ctx.JSON(http.StatusOK, list)
}

func (wh *WebHandler) toErrorMessage(err error) (int, *ErrorMessage) {

	var errorCode *service.ErrorCode
	errorMessage := &ErrorMessage{}

	if !errors.As(err, &errorCode) {
		errorMessage.Code = 100
		errorMessage.Message = err.Error()
		return http.StatusBadRequest, errorMessage
	}

	errorMessage.Code = errorCode.Code
	errorMessage.Message = errorCode.Message

	if len(errorCode.Args) > 0 {
		errorMessage.Message = fmt.Sprintf(errorCode.Message, errorCode.Args...)
	}

	return errorCode.Status, errorMessage
}

func (wh *WebHandler) beanToUser(user *inmemory.UserBean) User {
	return User{
		Id:       user.Id.String(),
		Username: user.Username,
	}
}

func (wh *WebHandler) beanToChannel(channel *inmemory.ChannelBean) Channel {
	members := make([]string, 0)
	for _, v := range channel.Members {
		members = append(members, v)
	}

	return Channel{
		Id:         channel.Id.String(),
		Members:    members,
		Name:       channel.Name,
		OwnerId:    channel.OwnerId,
		Visibility: channel.Visibility,
	}
}
