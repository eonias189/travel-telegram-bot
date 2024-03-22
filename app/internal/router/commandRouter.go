package router

import "github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"

func NewCommandRouter() *Router {
	return &Router{
		routes:     make(map[string]tgapi.HandlerFunc),
		onNotFound: func(ctx *tgapi.Context) error { return nil },
		getRoute: func(ctx *tgapi.Context) (string, bool) {
			if !ctx.Update.Message.IsCommand() {
				return "", false
			}
			return ctx.Update.Message.Command(), true
		},
	}
}
