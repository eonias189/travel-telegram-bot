package app

import "github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"

func handleNotes(opts AppHandlerOptions) {
	opts.CallbackRouter.Handle("notes", func(ctx *tgapi.Context) error {
		return ctx.SendString("пока не реализовано")
	})
}
