package applogger

import (
	"log/slog"
	"os"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/config"
)

func New(env config.LoggerEnv) *slog.Logger {
	var handler slog.Handler

	switch env {
	case config.LocalEnv:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case config.DevelopmentEnv:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	case config.ProductionEnv:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		handler = slog.Default().Handler()
	}

	return slog.New(handler)
}
