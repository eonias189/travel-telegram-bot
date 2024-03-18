package clearcontext

import (
	texthandler "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/textHandler"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/tgrouter"
)

func New() tgrouter.Middleware {
	return func(next tgrouter.HandlerFunc) tgrouter.HandlerFunc {
		return func(ctx *tgrouter.Context) error {
			if ctx.Update.Message.IsCommand() {
				texthandler.SetDialogContext(ctx, "")
			}
			return next(ctx)
		}
	}
}
