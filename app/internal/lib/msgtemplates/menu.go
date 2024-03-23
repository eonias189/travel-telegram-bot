package msgtempl

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func MenuMsg(chatId int64) tgbotapi.MessageConfig {
	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("путешествия", "trips"),
		),
	)
	msg := tgbotapi.NewMessage(chatId, "меню")
	msg.ReplyMarkup = btns
	return msg
}
