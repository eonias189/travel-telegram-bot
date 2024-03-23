package msgtempl

import (
	"fmt"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TripsMessage(senderId int64, trips []service.Trip) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(senderId, "путешествия")

	rows := make([][]tgbotapi.InlineKeyboardButton, len(trips)+2)
	for i, trip := range trips {
		rows[i+2] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(trip.Name, fmt.Sprintf("trip?id=%v", trip.Id)),
			tgbotapi.NewInlineKeyboardButtonData("удалить", fmt.Sprintf("delete-trip?id=%v", trip.Id)),
		)
	}

	rows[0] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "menu"))
	rows[1] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("новое путешествие", "new-trip"))

	btns := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyMarkup = btns
	return msg

}
