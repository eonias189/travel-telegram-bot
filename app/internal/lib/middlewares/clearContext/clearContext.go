package clearcontext

import "github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"

type DialogContexSetter interface {
	SetDialogContext(ctx *tgapi.Context, dialogContext string)
}

func NewAfterCleaner(dcs DialogContexSetter) tgapi.Middleware {
	return func(next tgapi.HandlerFunc) tgapi.HandlerFunc {
		return func(ctx *tgapi.Context) error {
			err := next(ctx)
			dcs.SetDialogContext(ctx, "")
			return err
		}
	}
}

func NewBeforeCleaner(dcs DialogContexSetter) tgapi.Middleware {
	return func(next tgapi.HandlerFunc) tgapi.HandlerFunc {
		return func(ctx *tgapi.Context) error {
			dcs.SetDialogContext(ctx, "")
			return next(ctx)
		}
	}
}
