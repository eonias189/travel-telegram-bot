package app

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/keyboards/menu"
	someimg "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/keyboards/someImg"
	clearcontext "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/clearContext"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/logger"
	texthandler "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/textHandler"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/tgrouter"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/view"
)

var (
	ErrInternal = errors.New("internal error")
)

type App struct {
	logger      *slog.Logger
	router      tgrouter.Router
	view        view.View
	textHandler *texthandler.TextHandler
}

func (a *App) Run(ctx context.Context, token string) error {
	a.logger.Info("starting")

	a.view.AddImage("menu", menu.New())
	a.view.AddImage("another", someimg.New())

	a.router.Use(logger.New(a.logger)).Use(clearcontext.New())

	a.router.On("start", func(ctx *tgrouter.Context) error {
		return ctx.SendString("starting")
	})
	a.router.On("reverse", func(ctx *tgrouter.Context) error {
		texthandler.SetDialogContext(ctx, "reverse")
		return nil
	})
	a.router.On("menu", func(ctx *tgrouter.Context) error {
		keyboard, ok := a.view.GetImage("menu")
		if !ok {
			return ErrInternal
		}
		return ctx.SetKeyboard(keyboard, "opening menu")
	})
	a.router.On("another", func(ctx *tgrouter.Context) error {
		keyboard, ok := a.view.GetImage("another")
		if !ok {
			return ErrInternal
		}
		return ctx.SetKeyboard(keyboard, "opening another")
	})
	a.router.On("close", func(ctx *tgrouter.Context) error {
		return ctx.CloseKeyboard("closing")
	})

	a.textHandler.OnContext("reverse", func(ctx *tgrouter.Context) error {
		// texthandler.SetDialogContext(ctx, "")
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

	return a.router.Run(ctx, token)
}

func (a *App) Close() {
	a.router.Close()
}

func New(l *slog.Logger) *App {
	texthandler := texthandler.New(&texthandler.Config{
		OnUnknownContext: func(ctx *tgrouter.Context) error {
			return ctx.SendString("non-context messages are not handling")
		},
	})
	router := tgrouter.NewRouter(&tgrouter.Config{
		OnText:  texthandler.ToHandler(),
		Workers: 5,
	})
	view := view.New()

	return &App{logger: l, router: router, textHandler: texthandler, view: view}
}
