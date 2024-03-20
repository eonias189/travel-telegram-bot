package tgapi

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx *Context) error

type Api struct {
	MiddlewareMixin
	onError   func(error)
	onCommand HandlerFunc
	onText    HandlerFunc
}

func (a *Api) OnError(handler func(error)) {
	if handler != nil {
		a.onError = handler
	}
}

func (a *Api) OnCommand(handler HandlerFunc) {
	if handler != nil {
		a.onCommand = handler
	}
}

func (a *Api) OnText(handler HandlerFunc) {
	if handler != nil {
		a.onText = handler
	}
}

func (a *Api) Run(ctx context.Context, token string) error {

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	updatesCh := bot.GetUpdatesChan(tgbotapi.UpdateConfig{})

	for {

		select {

		case <-ctx.Done():
			return ctx.Err()

		case update := <-updatesCh:

			if update.Message == nil {
				continue
			}

			c := &Context{
				Update: update,
				Bot:    bot,
				ctx:    context.Background(),
			}

			var err error

			if update.Message.IsCommand() {
				err = a.WithMiddlewares(a.onCommand)(c)
			} else {
				err = a.WithMiddlewares(a.onText)(c)
			}

			if err != nil {
				a.onError(err)
			}
		}
	}
}

func NewApi() *Api {
	return &Api{
		onError:   func(err error) {},
		onCommand: func(ctx *Context) error { return nil },
		onText:    func(ctx *Context) error { return nil },
	}
}
