package router

import (
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

// Route must match `route/arg`
func NewCallbackRouter() *Router {

	return &Router{
		routes:     make(map[string]tgapi.HandlerFunc),
		onNotFound: func(ctx *tgapi.Context) error { return nil },
		getRoute: func(ctx *tgapi.Context) (string, bool) {
			data := ctx.Update.CallbackData()
			split_route := strings.Split(data, "/")
			if len(split_route) == 0 {
				return "", false
			}
			if len(split_route) > 2 {
				return "", false
			}
			return split_route[0], true
		},
	}
}
