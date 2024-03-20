package logger

import (
	"context"
	"log/slog"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type DialogContextProvider interface {
	GetDialogContext(ctx *tgapi.Context) string
}

func New(logger *slog.Logger, dialogProvider DialogContextProvider) tgapi.Middleware {
	return func(next tgapi.HandlerFunc) tgapi.HandlerFunc {
		return func(ctx *tgapi.Context) error {
			logger.
				With(slog.String("dialog context", dialogProvider.GetDialogContext(ctx))).
				Log(context.TODO(), slog.LevelInfo, ctx.Update.Message.Text)
			return next(ctx)
		}
	}
}
