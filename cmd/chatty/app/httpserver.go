package app

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/knadh/koanf"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const (
	httpServerConfigKey = "http.server"
)

// HTTPServerConfig provides configuration for HTTP Server
type HTTPServerConfig struct {
	Address                string                  `json:"address"`
	SystemHTTPServerConfig *SystemHTTPServerConfig `json:"system,omitempty"`
	ReadTimeout            time.Duration           `json:"readTimeout"`
	WriteTimeout           time.Duration           `json:"writeTimeout"`
	IdleTimeout            time.Duration           `json:"idleTimeout"`
}

// SystemHTTPServerConfig provides configuration for the system HTTP server
type SystemHTTPServerConfig struct {
	Address string `json:"address"`
}

// ReadHTTPServerConfig reads HTTP server configuration
func ReadHTTPServerConfig(k *koanf.Koanf) (*HTTPServerConfig, error) {
	httpServerConfig := &HTTPServerConfig{}
	if err := k.UnmarshalWithConf(httpServerConfigKey, httpServerConfig, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, err
	}
	return httpServerConfig, nil
}

// ConfigureHTTPServers creates an HTTP server with standard middleware and a system HTTP server with health and metrics endpoints given an optional TLS listener
// returns the echo engine for serving API
func ConfigureHTTPServers(lifecycle fx.Lifecycle, shutdowner fx.Shutdowner, k *koanf.Koanf) (*echo.Echo, error) {
	httpConfig, err := ReadHTTPServerConfig(k)
	if err != nil {
		return nil, err
	}

	if httpConfig.SystemHTTPServerConfig != nil {
		systemEcho, err := newSystemHTTPServer(httpConfig.SystemHTTPServerConfig)
		if err != nil {
			return nil, err
		}

		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					if err := systemEcho.Start(systemEcho.Listener.Addr().String()); err != nil && err != http.ErrServerClosed {
						zap.L().Error("failed to start system HTTP server", zap.Error(err))
						if err := shutdowner.Shutdown(); err != nil {
							zap.L().Error("fx shutdown error", zap.Error(err))
						}
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return systemEcho.Shutdown(ctx)
			},
		})
	}

	e, err := newEcho(httpConfig)
	if err != nil {
		return nil, err
	}

	e.Use(middleware.Recover())

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := e.Start(e.Listener.Addr().String()); err != nil && err != http.ErrServerClosed {
					zap.L().Error("failed to start echo server", zap.Error(err))
					if err := shutdowner.Shutdown(); err != nil {
						zap.L().Error("fx shutdown error", zap.Error(err))
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})

	return e, nil
}

func newSystemHTTPServer(config *SystemHTTPServerConfig) (*echo.Echo, error) {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.GET("/health", getHealth)
	e.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))

	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	e.Listener = listener

	return e, nil
}

func getHealth(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func newEcho(config *HTTPServerConfig) (*echo.Echo, error) {
	listener, err := net.Listen("tcp", config.Address)
	if err != nil {
		return nil, err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Server.ReadTimeout = config.ReadTimeout
	e.Server.WriteTimeout = config.WriteTimeout
	e.Server.IdleTimeout = config.IdleTimeout

	e.Listener = listener

	return e, nil
}
