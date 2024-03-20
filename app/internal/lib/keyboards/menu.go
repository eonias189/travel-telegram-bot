package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Menu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/reverse"),
		tgbotapi.NewKeyboardButton("/another"),
	),
)
