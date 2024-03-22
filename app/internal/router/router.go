package router

import (
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type Router struct {
	tgapi.MiddlewareMixin
	routes     map[string]tgapi.HandlerFunc
	onNotFound tgapi.HandlerFunc
	getRoute   func(ctx *tgapi.Context) (string, bool)
}

func (r *Router) OnNotFound(handler tgapi.HandlerFunc) {
	if handler != nil {
		r.onNotFound = handler
	}
}

func (r *Router) Handle(route string, handler tgapi.HandlerFunc) {
	r.routes[route] = handler
}

func (r *Router) ToHandler() tgapi.HandlerFunc {
	return func(ctx *tgapi.Context) error {

		route, ok := r.getRoute(ctx)
		if !ok {
			return r.WithMiddlewares(r.onNotFound)(ctx)
		}

		handler, ok := r.routes[route]
		if !ok {
			return r.WithMiddlewares(r.onNotFound)(ctx)
		}

		return r.WithMiddlewares(handler)(ctx)
	}
}
