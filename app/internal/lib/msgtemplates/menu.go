package msgtempl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MenuButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("путешествия", "trips"),
		),
	)
}
