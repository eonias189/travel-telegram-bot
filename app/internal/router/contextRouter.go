package router

import (
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

type DialogContextProvider interface {
	GetDialogContext(ctx *tgapi.Context) string
}

func NewContextRouter(dcp DialogContextProvider) *Router {
	return &Router{
		routes:     make(map[string]tgapi.HandlerFunc),
		onNotFound: func(ctx *tgapi.Context) error { return nil },
		getRoute: func(ctx *tgapi.Context) (string, bool) {
			dialogContext := dcp.GetDialogContext(ctx)
			splitedDC := strings.Split(dialogContext, "?")

			if len(splitedDC) == 0 {
				return "", false
			}

			if len(splitedDC) > 2 {
				return "", false
			}

			return splitedDC[0], true
		},
	}
}
