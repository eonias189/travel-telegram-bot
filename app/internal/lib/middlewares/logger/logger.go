package logger

import (
	"context"
	"log/slog"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/tgrouter"
)

func log(l *slog.Logger, ctx *tgrouter.Context) {
	l.With(slog.String("message", ctx.Update.Message.Text)).
		With(slog.Int64("chat id", ctx.Update.Message.Chat.ID)).
		Log(context.TODO(), slog.LevelInfo, "handling request")
}

func New(l *slog.Logger) tgrouter.Middleware {
	return func(next tgrouter.HandlerFunc) tgrouter.HandlerFunc {
		return func(ctx *tgrouter.Context) error {
			log(l, ctx)
			return next(ctx)
		}
	}
}
