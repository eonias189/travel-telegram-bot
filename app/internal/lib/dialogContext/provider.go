package dialogcontext

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
)

var (
	ctxKey = "dialog context"
)

type StringCash interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type DialogContextProvider struct {
	cash StringCash
}

func (d *DialogContextProvider) Middleware() tgapi.Middleware {
	return func(next tgapi.HandlerFunc) tgapi.HandlerFunc {
		return func(ctx *tgapi.Context) error {
			id := ctx.Update.SentFrom().ID
			dialogContext, _ := d.cash.Get(fmt.Sprint(id))
			ctx = ctx.WithCtx(context.WithValue(ctx.Ctx(), ctxKey, dialogContext))
			return next(ctx)
		}
	}
}

func (d *DialogContextProvider) GetDialogContext(ctx *tgapi.Context) string {
	dialogCtx, _ := ctx.Ctx().Value(ctxKey).(string)
	return dialogCtx
}

func (d *DialogContextProvider) GetDialogContextQuery(ctx *tgapi.Context) url.Values {
	dctx, _ := ctx.Ctx().Value(ctxKey).(string)
	splitDCTX := strings.Split(dctx, "?")

	if len(splitDCTX) != 2 {
		return url.Values{}
	}

	query, err := url.ParseQuery(splitDCTX[1])
	if err != nil {
		return url.Values{}
	}

	return query
}

func (d *DialogContextProvider) SetDialogContext(ctx *tgapi.Context, dialogContext string) {
	id := ctx.Update.SentFrom().ID
	d.cash.Set(fmt.Sprint(id), dialogContext)
}

func NewProvider(cash StringCash) *DialogContextProvider {
	return &DialogContextProvider{
		cash: cash,
	}
}
