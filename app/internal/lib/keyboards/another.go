package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Another = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("/menu")),
)
