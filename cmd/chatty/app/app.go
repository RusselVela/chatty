package app

import (
	"embed"
	"github.com/RusselVela/chatty/internal/app/datasourcce/repository/inmemory"

	"github.com/RusselVela/chatty/internal/app/service"
	"github.com/RusselVela/chatty/internal/handler/web"

	"go.uber.org/fx"
)

//go:embed embedded/application.yaml
var applicationConfig embed.FS

// New initializes with the default config path
func New() *fx.App {
	return NewWithConfig("")
}

// NewWithConfig initializes with the specified config path
func NewWithConfig(cfgPath string) *fx.App {
	k, err := NewDefaultKoanf(&applicationConfig, cfgPath)
	if err != nil {
		panic(err)
	}

	app := fx.New(
		fx.Supply(k),
		fx.Provide(
			ConfigureLogger,
			web.GetSwagger,
			fx.Annotate(web.ConfigureHTTPServers, fx.As(new(web.EchoRouter))),
			fx.Annotate(service.NewChattyService, fx.As(new(web.ChattyService))),
			fx.Annotate(web.NewWebHandler, fx.As(new(web.ServerInterface))),
		),
		fx.Invoke(
			service.NewWebsocketHandler,
			service.SetupJWTSecret,
			inmemory.InitDatabase,
			web.RegisterHandlers,
			service.HandleMessages,
			LogAppStartStop,
		),
		fx.WithLogger(ConfigureFxLogger),
	)

	return app
}
