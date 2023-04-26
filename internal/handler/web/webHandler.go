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

func (wh *WebHandler) PublicPostWs(ctx echo.Context) error {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(*service.JWTCustomClaims)

	err := wh.service.HandleConnections(ctx, claims.Username)
	if err != nil {

	}
	return nil
}
