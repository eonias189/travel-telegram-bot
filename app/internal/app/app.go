package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/logger"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/tgrouter"
)

type App struct {
	l *slog.Logger
	r tgrouter.Router
}

func (a *App) Run(ctx context.Context, token string) error {
	a.l.Info("starting")

	a.r.Use(logger.New(a.l))
	a.r.Handle("start", func(ctx *tgrouter.Context) error {
		return ctx.SendString("starting")
	})
	a.r.Handle("err", func(ctx *tgrouter.Context) error {
		return fmt.Errorf("err")
	})

	return a.r.Run(ctx, token)
}

func (a *App) Close() {
	a.r.Close()
}

func New(l *slog.Logger) *App {
	router := tgrouter.NewRouter(&tgrouter.Config{
		OnText: func(ctx *tgrouter.Context) error {
			return ctx.SendString("пока не принимаю текстовые сообщения")
		},
		Workers: 5,
	})

	return &App{l: l, r: router}
}
