package tgrouter

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	Update     tgbotapi.Update
	Bot        *tgbotapi.BotAPI
	ctx        context.Context
	setContext func(context.Context)
}

func (c *Context) Ctx() context.Context {
	return c.ctx
}

func (c *Context) SetContextValue(key string, value any) {
	c.setContext(context.WithValue(c.ctx, key, value))
}

func (c *Context) SendString(s string) error {
	msg := tgbotapi.NewMessage(c.Update.Message.Chat.ID, s)
	_, err := c.Bot.Send(msg)
	return err
}
