package contextrouter

import (
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type DialogContextProvider interface {
	GetDialogContext(ctx *tgapi.Context) string
}

type ContextRouter struct {
	tgapi.MiddlewareMixin
	routes         map[string]tgapi.HandlerFunc
	dialogProvider DialogContextProvider
	onNotFound     tgapi.HandlerFunc
}

func (cr *ContextRouter) Handle(ctx string, handler tgapi.HandlerFunc) {
	if handler != nil {
		cr.routes[ctx] = handler
	}
}

func (cr *ContextRouter) OnNotFound(handler tgapi.HandlerFunc) {
	if handler != nil {
		cr.onNotFound = handler
	}
}

func (cr *ContextRouter) ToHandler() tgapi.HandlerFunc {
	return func(ctx *tgapi.Context) error {
		c := cr.dialogProvider.GetDialogContext(ctx)

		handler, ok := cr.routes[c]
		if !ok {
			return cr.WithMiddlewares(cr.onNotFound)(ctx)
		}

		return cr.WithMiddlewares(handler)(ctx)
	}
}

func New(dialogContextProvider DialogContextProvider) *ContextRouter {
	return &ContextRouter{
		routes:         make(map[string]tgapi.HandlerFunc),
		dialogProvider: dialogContextProvider,
		onNotFound:     func(ctx *tgapi.Context) error { return nil },
	}
}
