package commandrouter

import (
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type CommandRouter struct {
	tgapi.MiddlewareMixin
	routes     map[string]tgapi.HandlerFunc
	onNotFound tgapi.HandlerFunc
}

func (cr *CommandRouter) OnNotFound(handler tgapi.HandlerFunc) {
	if handler != nil {
		cr.onNotFound = handler
	}
}

func (cr *CommandRouter) Handle(command string, handler tgapi.HandlerFunc) {
	if handler != nil {
		cr.routes[command] = handler
	}
}

func (cr *CommandRouter) ToHandler() tgapi.HandlerFunc {
	return func(ctx *tgapi.Context) error {
		command := ctx.Update.Message.Command()

		handler, ok := cr.routes[command]
		if !ok {
			return cr.WithMiddlewares(cr.onNotFound)(ctx)
		}

		return cr.WithMiddlewares(handler)(ctx)
	}
}

func New() *CommandRouter {
	return &CommandRouter{
		routes:     make(map[string]tgapi.HandlerFunc),
		onNotFound: func(ctx *tgapi.Context) error { return nil },
	}
}
