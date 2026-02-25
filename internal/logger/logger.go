package logger

import (
	"log/slog"
	"os"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log/global"
)

type LogConfig struct {
	Environment string
	LogLevel    string
}

func toSlogLevel(logLevel string) slog.Leveler {
	switch logLevel {
	case "ERROR":
		return slog.LevelError
	case "WARN":
		return slog.LevelWarn
	case "INFO":
		return slog.LevelInfo
	case "DEBUG":
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

func New(cfg LogConfig) *slog.Logger {
	option := &slog.HandlerOptions{
		// AddSource: true,
		Level: toSlogLevel(cfg.LogLevel),
	}

	if cfg.Environment != "production" {
		return slog.New(slog.NewTextHandler(os.Stdout, option))
	}

	loggerProvider := global.GetLoggerProvider()

	return slog.New(otelslog.NewHandler("mall", otelslog.WithLoggerProvider(loggerProvider)))
}
