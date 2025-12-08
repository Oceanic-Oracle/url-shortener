package logger

import (
	"log/slog"
	"os"

	"shortener/internal/logctx"
)

const (
	envDebug = "debug"
	envInfo  = "info"
	envError = "error"
)

func SetupLogger(env string) *slog.Logger {
	var level slog.Level

	switch env {
	case envDebug:
		level = slog.LevelDebug
	case envInfo:
		level = slog.LevelInfo
	case envError:
		level = slog.LevelError
	default:
		level = slog.LevelDebug
	}

	return slog.New(logctx.NewHandlerMiddleware(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}),
	))
}
