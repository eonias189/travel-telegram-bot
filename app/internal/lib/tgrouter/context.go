package tgrouter

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	Update  tgbotapi.Update
	Bot     *tgbotapi.BotAPI
	Context context.Context
}

func (ctx *Context) SendString(s string) error {
	msg := tgbotapi.NewMessage(ctx.Update.Message.Chat.ID, s)
	_, err := ctx.Bot.Send(msg)
	return err
}
