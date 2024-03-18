package someimg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func New() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/menu")),
	)
}
