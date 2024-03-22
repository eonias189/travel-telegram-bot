package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	dialogcontext "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/dialogContext"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/keyboards"
	clearcontext "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/clearContext"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/logger"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/router"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgapi"
	"github.com/redis/go-redis/v9"
)

var (
	ErrInternal = errors.New("internal error")
)

type App struct {
	logger *slog.Logger
	api    *tgapi.Api
	rdb    *redis.Client
}

func (a *App) handleAll() {
	cash := service.NewRedisCash(a.rdb.Conn(), time.Hour)
	dialogContextProvider := dialogcontext.NewProvider(cash)

	cmdr := router.NewCommandRouter()
	cmdr.Use(clearcontext.NewBeforeCleaner(dialogContextProvider))
	cmdr.OnNotFound(func(ctx *tgapi.Context) error {
		return ctx.SendString(fmt.Sprintf("комманда %v не найдена", ctx.Update.Message.Command()))
	})

	ctxr := router.NewContextRouter(dialogContextProvider)
	ctxr.OnNotFound(func(ctx *tgapi.Context) error {
		return ctx.SendString("сообщения вне контекста не обрабатываются")
	})

	cbr := router.NewCallbackRouter()
	cbr.Use(clearcontext.NewBeforeCleaner(dialogContextProvider))

	a.api.OnCommand(cmdr.ToHandler())
	a.api.OnCallback(cbr.ToHandler())
	a.api.OnText(ctxr.ToHandler())

	a.api.Use(logger.New(a.logger, dialogContextProvider))
	a.api.Use(dialogContextProvider.Middleware())

	cmdr.Handle("start", func(ctx *tgapi.Context) error {
		return ctx.SendString("starting")
	})

	cmdr.Handle("reverse", func(ctx *tgapi.Context) error {
		dialogContextProvider.SetDialogContext(ctx, "reverse")
		return nil
	})

	cmdr.Handle("menu", func(ctx *tgapi.Context) error {
		return ctx.SendWithInlineKeyboard("opening menu", keyboards.Menu)
	})

	cbr.Handle("print", func(ctx *tgapi.Context) error {
		arg := ctx.CallbackArg()
		return ctx.SendString(arg)
	})

	cmdr.Handle("close", func(ctx *tgapi.Context) error {
		return ctx.CloseKeyboard("closing")
	})

	ctxr.Handle("reverse", func(ctx *tgapi.Context) error {
		text := ctx.Update.Message.Text

		reverse := func(s string) string {
			var res string
			for i := len(s) - 1; i >= 0; i-- {
				res += string(s[i])
			}
			return res
		}

		return ctx.SendString(reverse(text))

	})
}

func (a *App) Run(ctx context.Context, token string) error {
	a.logger.Info("starting")

	a.handleAll()
	a.api.Run(ctx, token)
	return nil
}

func New(rdb *redis.Client, logger *slog.Logger) *App {
	return &App{logger: logger, rdb: rdb, api: tgapi.NewApi()}
}
