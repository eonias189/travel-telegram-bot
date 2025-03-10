package tgapi

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerFunc func(ctx *Context) error

type Api struct {
	MiddlewareMixin
	onError    func(ctx *Context, err error)
	onCommand  HandlerFunc
	onCallback HandlerFunc
	onText     HandlerFunc
}

func (a *Api) OnError(handler func(ctx *Context, err error)) {
	if handler != nil {
		a.onError = handler
	}
}

func (a *Api) OnCommand(handler HandlerFunc) {
	if handler != nil {
		a.onCommand = handler
	}
}

func (a *Api) OnCallback(handler HandlerFunc) {
	if handler != nil {
		a.onCallback = handler
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

			c := &Context{
				Update: update,
				Bot:    bot,
				ctx:    context.Background(),
			}
			var err error

			if cbData := update.CallbackData(); cbData != "" {
				err = a.WithMiddlewares(a.onCallback)(c)
			} else if update.Message == nil {
				return nil
			} else if update.Message.IsCommand() {
				err = a.WithMiddlewares(a.onCommand)(c)
			} else {
				err = a.WithMiddlewares(a.onText)(c)
			}

			if err != nil {
				a.onError(c, err)
			}
		}
	}
}

func NewApi() *Api {
	return &Api{
		onError:    func(ctx *Context, err error) {},
		onCommand:  func(ctx *Context) error { return nil },
		onCallback: func(ctx *Context) error { return nil },
		onText:     func(ctx *Context) error { return nil },
	}
}
