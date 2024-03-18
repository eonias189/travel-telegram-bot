package texthandler

import (
	"sync"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/tgrouter"
)

var (
	contextKey = "dialog context"
)

func SetDialogContext(ctx *tgrouter.Context, c string) {
	ctx.SetContextValue(contextKey, c)
}

func GetDialogContext(ctx *tgrouter.Context) string {
	c, _ := ctx.Ctx().Value(contextKey).(string)
	return c
}

type TextHandler struct {
	mu       sync.RWMutex
	contexts map[string]tgrouter.HandlerFunc
	cfg      *Config
}

func (t *TextHandler) OnContext(ctx string, handler tgrouter.HandlerFunc) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	t.contexts[ctx] = handler
}

func (t *TextHandler) ToHandler() tgrouter.HandlerFunc {
	return func(ctx *tgrouter.Context) error {
		c := GetDialogContext(ctx)

		handler, ok := t.contexts[c]
		if !ok {
			return t.cfg.OnUnknownContext(ctx)
		}

		return handler(ctx)
	}
}

func New(cfg *Config) *TextHandler {
	return &TextHandler{
		contexts: make(map[string]tgrouter.HandlerFunc),
		cfg:      handleConfig(cfg),
	}
}
