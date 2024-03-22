package router

import "github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"

type DialogContextProvider interface {
	GetDialogContext(ctx *tgapi.Context) string
}

func NewContextRouter(dcp DialogContextProvider) *Router {
	return &Router{
		routes:     make(map[string]tgapi.HandlerFunc),
		onNotFound: func(ctx *tgapi.Context) error { return nil },
		getRoute: func(ctx *tgapi.Context) (string, bool) {
			return dcp.GetDialogContext(ctx), true
		},
	}
}
