package tgapi

import (
	"context"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Context struct {
	Update tgbotapi.Update
	Bot    *tgbotapi.BotAPI
	ctx    context.Context
}

func (c *Context) Ctx() context.Context {
	return c.ctx
}

func (c *Context) WithCtx(ctx context.Context) *Context {
	return &Context{
		Update: c.Update,
		Bot:    c.Bot,
		ctx:    ctx,
	}
}

func (c *Context) SendString(s string) error {
	msg := tgbotapi.NewMessage(c.Update.SentFrom().ID, s)
	_, err := c.Bot.Send(msg)
	return err
}

func (c *Context) SendWithInlineKeyboard(text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(c.Update.SentFrom().ID, text)
	msg.ReplyMarkup = keyboard
	_, err := c.Bot.Send(msg)
	return err
}

func (c *Context) CallbackArg() string {
	data := c.Update.CallbackData()
	splitData := strings.Split(data, "/")
	if len(splitData) != 2 {
		return ""
	}
	return splitData[1]
}
