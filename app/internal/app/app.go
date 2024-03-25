package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	dialogcontext "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/dialogContext"
	clearcontext "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/clearContext"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/lib/middlewares/logger"
	msgtempl "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/msgtemplates"
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

func NewId() int64 {
	return time.Now().UnixNano()
}

func (a *App) handleAll() {
	cash := service.NewRedisCash(a.rdb.Conn(), service.CashOptions{Prefix: "dialog-context", ExpirationTime: time.Hour})
	dialogContextProvider := dialogcontext.NewProvider(cash)
	userService := service.NewUserServive(a.rdb)
	tripService := service.NewTripService(a.rdb)
	locationService := service.NewLocationService(a.rdb)

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

	handlerOptions := AppHandlerOptions{
		CommandRouter:  cmdr,
		ContextRouter:  ctxr,
		CallbackRouter: cbr,
		Dcs:            dialogContextProvider,
		Dcqp:           dialogContextProvider,
		Logger:         a.logger,
	}

	a.api.OnCommand(cmdr.ToHandler())
	a.api.OnCallback(cbr.ToHandler())
	a.api.OnText(ctxr.ToHandler())
	a.api.OnError(func(ctx *tgapi.Context, err error) {
		a.logger.Error(err.Error())
		ctx.SendString("произошла ошибка")
	})

	a.api.Use(logger.New(a.logger, dialogContextProvider))
	a.api.Use(dialogContextProvider.Middleware())

	cmdr.Handle("start", func(ctx *tgapi.Context) error {
		return ctx.SendMessage(msgtempl.MenuMsg(ctx.SenderID()))
	})

	cmdr.Handle("menu", func(ctx *tgapi.Context) error {
		return ctx.SendMessage(msgtempl.MenuMsg(ctx.SenderID()))
	})

	cbr.Handle("menu", func(ctx *tgapi.Context) error {
		return ctx.SendMessage(msgtempl.MenuMsg(ctx.SenderID()))
	})

	cmdr.Handle("otmena", func(ctx *tgapi.Context) error {
		dialogContextProvider.SetDialogContext(ctx, "")
		err := ctx.SendString("отмена контекста")
		if err != nil {
			return err
		}
		return ctx.SendMessage(msgtempl.MenuMsg(ctx.SenderID()))
	})

	handleProfile(handlerOptions, userService)
	handleTrips(handlerOptions, userService, tripService)
	handleLocations(handlerOptions, locationService)
	handleNotes(handlerOptions)
	handleFriends(handlerOptions, userService, tripService)
	handleGeoTools(handlerOptions, locationService)
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
