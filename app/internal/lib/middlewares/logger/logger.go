package logger

import (
	"log/slog"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type DialogContextProvider interface {
	GetDialogContext(ctx *tgapi.Context) string
}

func New(logger *slog.Logger, dialogProvider DialogContextProvider) tgapi.Middleware {
	return func(next tgapi.HandlerFunc) tgapi.HandlerFunc {
		return func(ctx *tgapi.Context) error {

			attrs := []any{
				slog.String("dialog_context", dialogProvider.GetDialogContext(ctx)),
			}

			if data := ctx.Update.CallbackData(); data != "" {
				attrs = append(attrs, slog.String("callback_data", data))
			} else {
				attrs = append(attrs, slog.String("message", ctx.Update.Message.Text))
			}

			logger.Info("handling message", attrs...)
			return next(ctx)
		}
	}
}
