package logger

import (
	"context"
	"log/slog"

	texthandler "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/textHandler"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgrouter"
)

func log(l *slog.Logger, ctx *tgrouter.Context) {
	l.
		With(slog.String("dialog context", texthandler.GetDialogContext(ctx))).
		Log(context.TODO(), slog.LevelInfo, ctx.Update.Message.Text)
}

func New(l *slog.Logger) tgrouter.Middleware {
	return func(next tgrouter.HandlerFunc) tgrouter.HandlerFunc {
		return func(ctx *tgrouter.Context) error {
			log(l, ctx)
			return next(ctx)
		}
	}
}
