package tgapi

import (
	"context"
	"net/url"
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

func (c *Context) SenderID() int64 {
	return c.Update.SentFrom().ID
}

func (c *Context) SendMessage(msg tgbotapi.MessageConfig) error {
	_, err := c.Bot.Send(msg)
	return err
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

func (c *Context) CallbackQuery() url.Values {
	data := c.Update.CallbackData()
	splitData := strings.Split(data, "?")
	if len(splitData) != 2 {
		return url.Values{}
	}

	query, err := url.ParseQuery(splitData[1])
	if err != nil {
		return url.Values{}
	}
	return query
}
