package app

import (
	"context"

	"github.com/knadh/koanf"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const loggingConfigKey = "logging"
const fxLoggingConfigKey = "fx.logging"

// Config contains config properties for logging
type Config struct {
	Level         string `json:"level"`
	UseJSONFormat bool   `json:"useJsonFormat"`
}

// FxLoggingConfig contains config properties for Fx
type FxLoggingConfig struct {
	Level string `json:"level"`
}

// ConfigureLogger reads logging configuration and initializes logger
func ConfigureLogger(k *koanf.Koanf) (*zap.Logger, error) {
	loggingConfig, err := readConfig(k)
	if err != nil {
		return nil, err
	}
	return initZap(loggingConfig)
}

func readConfig(k *koanf.Koanf) (*Config, error) {
	loggingConfig := &Config{}
	if err := k.UnmarshalWithConf(loggingConfigKey, loggingConfig, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return nil, err
	}
	return loggingConfig, nil
}

func initZap(c *Config) (*zap.Logger, error) {
	logger, err := newZapLogger(c)
	if err != nil {
		return nil, err
	}
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(zap.L())
	return zap.L(), nil
}

func newZapLogger(c *Config) (*zap.Logger, error) {
	l, err := zap.ParseAtomicLevel(c.Level)
	if err != nil {
		return nil, err
	}
	zapConfig := zap.Config{
		Level:            l,
		Development:      false,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if c.UseJSONFormat {
		zapConfig.Encoding = "json"
		zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	} else {
		zapConfig.Encoding = "console"
		zapConfig.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		zapConfig.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	}
	return zapConfig.Build()
}

// LogAppStartStop logs a message to the configured logger when the applications starts and stops
func LogAppStartStop(lifecyle fx.Lifecycle, l *zap.Logger, k *koanf.Koanf) {
	lifecyle.Append(fx.Hook{
		OnStart: func(c context.Context) error {
			l.Info(k.String("service.name") + " started")
			return nil
		},
		OnStop: func(c context.Context) error {
			l.Info(k.String("service.name") + " stopped")
			return nil
		},
	})
}

// ConfigureFxLogger reads configuration and initializes a Fx logger
func ConfigureFxLogger(k *koanf.Koanf) fxevent.Logger {
	loggingConfig, err := readConfig(k)
	if err != nil {
		return &fxevent.NopLogger
	}
	fxLoggingConfig := &FxLoggingConfig{}
	if err := k.UnmarshalWithConf(fxLoggingConfigKey, fxLoggingConfig, koanf.UnmarshalConf{Tag: "json"}); err != nil {
		return &fxevent.NopLogger
	}

	loggingConfig.Level = fxLoggingConfig.Level
	l, err := newZapLogger(loggingConfig)
	if err != nil {
		return &fxevent.NopLogger
	}

	return &fxevent.ZapLogger{Logger: l.Named("fx")}
}
