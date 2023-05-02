// Package web provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.2 DO NOT EDIT.
package web

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	JwtAuthScopes = "jwtAuth.Scopes"
)

// Channel defines model for Channel.
type Channel struct {
	Id         string   `json:"id"`
	Members    []string `json:"members"`
	Name       string   `json:"name"`
	OwnerId    string   `json:"ownerId"`
	Visibility string   `json:"visibility"`
}

// ChannelCreation defines model for ChannelCreation.
type ChannelCreation struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ChannelSubscription defines model for ChannelSubscription.
type ChannelSubscription = map[string]interface{}

// ErrorMessage defines model for ErrorMessage.
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// LoginSchema defines model for LoginSchema.
type LoginSchema struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// SignUpSchema defines model for SignUpSchema.
type SignUpSchema = LoginSchema

// SuccessChannelCreation defines model for SuccessChannelCreation.
type SuccessChannelCreation struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	OwnerId string `json:"ownerId"`
}

// SuccessChannelSubscription defines model for SuccessChannelSubscription.
type SuccessChannelSubscription struct {
	Name string `json:"name"`
	Ok   bool   `json:"ok"`
}

// SuccessGetChannels defines model for SuccessGetChannels.
type SuccessGetChannels struct {
	Channels []Channel `json:"channels"`
}

// SuccessGetUsers defines model for SuccessGetUsers.
type SuccessGetUsers struct {
	Users []User `json:"users"`
}

// SuccessLoginSchema defines model for SuccessLoginSchema.
type SuccessLoginSchema struct {
	Token string `json:"token"`
}

// SuccessSignupSchema defines model for SuccessSignupSchema.
type SuccessSignupSchema struct {
	Id       string `json:"id"`
	Ok       *bool  `json:"ok,omitempty"`
	Username string `json:"username"`
}

// SuccessWsConnectionSchema defines model for SuccessWsConnectionSchema.
type SuccessWsConnectionSchema struct {
	Ok        bool   `json:"ok"`
	Timestamp string `json:"timestamp"`
}

// UnauthorizedSchema defines model for UnauthorizedSchema.
type UnauthorizedSchema struct {
	Message string `json:"message"`
}

// User defines model for User.
type User struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

// N200SuccessChannelSubscribe defines model for 200SuccessChannelSubscribe.
type N200SuccessChannelSubscribe = SuccessChannelSubscription

// N200SuccessGetChannels defines model for 200SuccessGetChannels.
type N200SuccessGetChannels = SuccessGetChannels

// N200SuccessGetUsers defines model for 200SuccessGetUsers.
type N200SuccessGetUsers = SuccessGetUsers

// N200SuccessLogin defines model for 200SuccessLogin.
type N200SuccessLogin = SuccessLoginSchema

// N200SuccessWsConnection defines model for 200SuccessWsConnection.
type N200SuccessWsConnection = SuccessWsConnectionSchema

// N200SuccessWsToken defines model for 200SuccessWsToken.
type N200SuccessWsToken = SuccessLoginSchema

// N201SuccessChannelCreation defines model for 201SuccessChannelCreation.
type N201SuccessChannelCreation = SuccessChannelCreation

// N201SuccessSignUp defines model for 201SuccessSignUp.
type N201SuccessSignUp = SuccessSignupSchema

// N400BadRequest defines model for 400BadRequest.
type N400BadRequest = ErrorMessage

// N401UnauthorizedError defines model for 401UnauthorizedError.
type N401UnauthorizedError = UnauthorizedSchema

// N500InternalServerError defines model for 500InternalServerError.
type N500InternalServerError = ErrorMessage

// ChannelCreationRequest defines model for ChannelCreationRequest.
type ChannelCreationRequest = ChannelCreation

// ChannelSubscriptionRequest defines model for ChannelSubscriptionRequest.
type ChannelSubscriptionRequest = ChannelSubscription

// LoginRequest defines model for LoginRequest.
type LoginRequest = LoginSchema

// SignUpRequest defines model for SignUpRequest.
type SignUpRequest = SignUpSchema

// PublicPostChannelsJSONRequestBody defines body for PublicPostChannels for application/json ContentType.
type PublicPostChannelsJSONRequestBody = ChannelCreation

// PublicPostChannelsSubscribeJSONRequestBody defines body for PublicPostChannelsSubscribe for application/json ContentType.
type PublicPostChannelsSubscribeJSONRequestBody = ChannelSubscription

// PublicPostSignupJSONRequestBody defines body for PublicPostSignup for application/json ContentType.
type PublicPostSignupJSONRequestBody = SignUpSchema

// PublicPostTokenJSONRequestBody defines body for PublicPostToken for application/json ContentType.
type PublicPostTokenJSONRequestBody = LoginSchema

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Returns a list of channels
	// (GET /v1/chatty/channels)
	PublicGetChannels(ctx echo.Context) error
	// Creates a new channel
	// (POST /v1/chatty/channels)
	PublicPostChannels(ctx echo.Context) error
	// Subscribe to a channel
	// (POST /v1/chatty/channels/{id}/subscribe)
	PublicPostChannelsSubscribe(ctx echo.Context, id string) error
	// Registers a new user
	// (POST /v1/chatty/signup)
	PublicPostSignup(ctx echo.Context) error
	// Returns a token
	// (POST /v1/chatty/token)
	PublicPostToken(ctx echo.Context) error
	// Returns a list of users
	// (GET /v1/chatty/users)
	PublicGetUsers(ctx echo.Context) error
	// Connects to Chatty to send and receive messages
	// (GET /v1/chatty/ws)
	PublicGetWs(ctx echo.Context) error
	// Returns a token for authentication with the Webscoket endpoint
	// (GET /v1/chatty/ws/token)
	PublicGetWsToken(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PublicGetChannels converts echo context to params.
func (w *ServerInterfaceWrapper) PublicGetChannels(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicGetChannels(ctx)
	return err
}

// PublicPostChannels converts echo context to params.
func (w *ServerInterfaceWrapper) PublicPostChannels(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicPostChannels(ctx)
	return err
}

// PublicPostChannelsSubscribe converts echo context to params.
func (w *ServerInterfaceWrapper) PublicPostChannelsSubscribe(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicPostChannelsSubscribe(ctx, id)
	return err
}

// PublicPostSignup converts echo context to params.
func (w *ServerInterfaceWrapper) PublicPostSignup(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicPostSignup(ctx)
	return err
}

// PublicPostToken converts echo context to params.
func (w *ServerInterfaceWrapper) PublicPostToken(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicPostToken(ctx)
	return err
}

// PublicGetUsers converts echo context to params.
func (w *ServerInterfaceWrapper) PublicGetUsers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicGetUsers(ctx)
	return err
}

// PublicGetWs converts echo context to params.
func (w *ServerInterfaceWrapper) PublicGetWs(ctx echo.Context) error {
	var err error

	ctx.Set(JwtAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicGetWs(ctx)
	return err
}

// PublicGetWsToken converts echo context to params.
func (w *ServerInterfaceWrapper) PublicGetWsToken(ctx echo.Context) error {
	var err error

	ctx.Set(JwtAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicGetWsToken(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/chatty/channels", wrapper.PublicGetChannels)
	router.POST(baseURL+"/v1/chatty/channels", wrapper.PublicPostChannels)
	router.POST(baseURL+"/v1/chatty/channels/:id/subscribe", wrapper.PublicPostChannelsSubscribe)
	router.POST(baseURL+"/v1/chatty/signup", wrapper.PublicPostSignup)
	router.POST(baseURL+"/v1/chatty/token", wrapper.PublicPostToken)
	router.GET(baseURL+"/v1/chatty/users", wrapper.PublicGetUsers)
	router.GET(baseURL+"/v1/chatty/ws", wrapper.PublicGetWs)
	router.GET(baseURL+"/v1/chatty/ws/token", wrapper.PublicGetWsToken)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
