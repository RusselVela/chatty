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
	Name string `json:"name"`
	Ok   bool   `json:"ok"`
}

// SuccessChannelSubscription defines model for SuccessChannelSubscription.
type SuccessChannelSubscription struct {
	Name string `json:"name"`
	Ok   bool   `json:"ok"`
}

// SuccessLoginSchema defines model for SuccessLoginSchema.
type SuccessLoginSchema struct {
	Ok    bool   `json:"ok"`
	Token string `json:"token"`
}

// SuccessSignupSchema defines model for SuccessSignupSchema.
type SuccessSignupSchema struct {
	Id       string `json:"id"`
	Ok       bool   `json:"ok"`
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

// N200SuccessChannelSubscribe defines model for 200SuccessChannelSubscribe.
type N200SuccessChannelSubscribe = SuccessChannelSubscription

// N200SuccessfulLogin defines model for 200SuccessfulLogin.
type N200SuccessfulLogin = SuccessLoginSchema

// N200SuccessfulWsConnection defines model for 200SuccessfulWsConnection.
type N200SuccessfulWsConnection = SuccessWsConnectionSchema

// N201SuccessChannelCreation defines model for 201SuccessChannelCreation.
type N201SuccessChannelCreation = SuccessChannelCreation

// N201SuccessfulSignUp defines model for 201SuccessfulSignUp.
type N201SuccessfulSignUp = SuccessSignupSchema

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
	// Creates a new channel
	// (POST /v1/chatty/channels)
	PublicPostChannels(ctx echo.Context) error
	// Subscribe to a channel
	// (POST /v1/chatty/channels/{name}/subscribe)
	PublicPostChannelsSubscribe(ctx echo.Context, name string) error
	// Connects to Chatty to send and receive messages
	// (GET /v1/chatty/chats/ws)
	PublicGetWs(ctx echo.Context) error
	// Registers a new user
	// (POST /v1/chatty/signup)
	PublicPostSignup(ctx echo.Context) error
	// Returns a token
	// (POST /v1/chatty/token)
	PublicPostToken(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
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
	// ------------- Path parameter "name" -------------
	var name string

	err = runtime.BindStyledParameterWithLocation("simple", false, "name", runtime.ParamLocationPath, ctx.Param("name"), &name)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter name: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PublicPostChannelsSubscribe(ctx, name)
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

	router.POST(baseURL+"/v1/chatty/channels", wrapper.PublicPostChannels)
	router.POST(baseURL+"/v1/chatty/channels/:name/subscribe", wrapper.PublicPostChannelsSubscribe)
	router.GET(baseURL+"/v1/chatty/chats/ws", wrapper.PublicGetWs)
	router.POST(baseURL+"/v1/chatty/signup", wrapper.PublicPostSignup)
	router.POST(baseURL+"/v1/chatty/token", wrapper.PublicPostToken)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xYzW7jNhB+FYIt0IsROe32olM3QVtksUWDONscAh9oaSwxkUgth3LgBn73gkPZkizJ",
	"VrwOujebHA4/zc83M3zlkc4LrUBZ5OErN/C1BLRXOpZAC9epUAqyawPCSq3u/L7bibSyoOinKIpMRiQQ",
	"PKFWbg2jFHLhfv1oYMlD/kNQXxX4XQz21PPNZjPhMWBkZEELIb9PgVWwmNUscqLABFPwwiJ/nG8mW6Cz",
	"crE7/E5gm1eMAIxefAHuTwf2Z53IswMlpTMvehxg5qQdlplM1Jfi3GC81tFoUCaqLLiTM4CFVugD8efp",
	"dFZGESC23bCA8yHt0z/o5t39MYVl7dIa6bLMyBXnRnjEv7TNJDLcwejgesBrrRRE/sh54TV1D6P8CxBF",
	"AgxB2QbSbO2xXra9sWOI93H2IQKqRCrmidvwlmXmA/zcwGaUB8Pm8/tdL3+YTq9EfO4k/t0YbSqP9aG5",
	"EjHb3kkgLr8oUdpUG/kvxHT6bFiamoft8+nhnln9DJQIuUSUKmHaMKlWIpPkxV+n0xtlwSiRzcCswJwX",
	"5zGbbe9m/nLmb3dylYae+uuWCqMLMLYq0ErkxIB2XQAPOVojVeK+zi90NohXv5bSQMzDR3+8Ep5vD3G9",
	"eILIDpTVhtJarvWxHZCRjmk1hqUoM8vD6e4qqSwkYJyOvD5+GDSpq+X7cDcpsgOnEIgv2sTu91KbXFge",
	"1ouTri1LdJ7KR0DbSU5qhX34WlXRhVmW/b3k4eMbCvvcqRkkyZFRop8bywutMxBqKEj0c/+nDJfN/wPH",
	"Qc/3XzPhxBQjw2EPlH7m2/MHULXovANLxm8wylvCkcA1YlLGh0D2FO7xFpQ5oBV5MRJULd8HqIfjO0hG",
	"08UwTziuhag00q7pFq/46cV+LG3qfi5AGDB/bIPi08M9r+iZPp926whJrS0800u11NtCIiIqJJALmfGQ",
	"ryATF6ZEhOy3xK1dRDrnkyo1+B3tsH8gE7xTMj7e3lCzqfO8VK4qAZOKGRAZc/ZkL9KmrAAwSBa22Z7G",
	"n5Bdp8LaNUMwKxm5kFiBQa99dUlBV4ASheQh/+ViejElIrMpGSZYXQYRnQ+qftfTqfadxl7bRO0S7o08",
	"pN8QS93EPOS35SKT0a1Ge73VOGlMoeuhItsaVIOBKbU7Q1wO66vkguHW03dXxzW0WzDfaBw/NdCNUIyW",
	"eS7M+oBRrUjQhbo3J5+7Uz3eCl5dlG0CbM5OhTAiB+uCxtWf7nTmzjC9ZDaF7ZWtwdbxiiLKtGkdyBXj",
	"1IloTQmTRu+0n7TzyUAkzeoJOgXm2Mx1cfrZdXQ2lch2IeVgNVFiAZFcSojZYs38t4+IwFnjw04Nxb53",
	"iJ6Rdkw4Ds6930E8ztqPG28JSIvBC5khgT728IUInV5PWQOO+xPsg+eM0y27Pw+fbNgPYwimdzg6h1eq",
	"SkZZvKthj3OXWA0G6RiWchlUzISKmYEI5ApYVTPxuCurB5vBOnAHiUTHLhVpufw9kIS+UTol89rPV9/I",
	"/fVc/x1k2YAJjzhm19MO+cWWRjmVfky2mi2IXWPXUyxLY1MwNbXiAZ/d000nuKz1+rn55iT+vH3JHJmH",
	"7+iwlm17fUUZ6xT4slsa1xxmOhJZqtEGrhmb787tO68gNezj7Q3WFbfSvZlv/gsAAP//jfVjoNcXAAA=",
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
