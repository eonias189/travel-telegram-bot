package msgtempl

import (
	"fmt"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TripsMessage(senderId int64, trips []service.Trip) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(senderId, "путешествия")

	rows := make([][]tgbotapi.InlineKeyboardButton, len(trips)+3)
	for i, trip := range trips {
		rows[i+3] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(trip.Name, fmt.Sprintf("trip?id=%v", trip.Id)))
	}

	rows[0] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "menu"))
	rows[1] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("присоединиться к путешествию", "join-to-trip"))
	rows[2] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("новое путешествие", "new-trip"))

	btns := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg.ReplyMarkup = btns
	return msg

}

func TripMessage(senderId int64, trip service.Trip) tgbotapi.MessageConfig {
	text := fmt.Sprintf(`%v
Описание: %v`, trip.Name, trip.Description)

	rows := [][]tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "trips")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("локации", fmt.Sprintf(`locations?tripId=%v`, trip.Id))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("заметки", fmt.Sprintf(`notes?tripId=%v`, trip.Id))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("проложить маршрут", fmt.Sprintf("get-route?tripId=%v", trip.Id))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("пригласить друзей", fmt.Sprintf(`share-trip?tripId=%v`, trip.Id))),
	}

	if trip.Creator == senderId {
		rows = append(rows,
			tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("удалить", fmt.Sprintf("delete-trip?id=%v", trip.Id))),
		)
	}

	msg := tgbotapi.NewMessage(senderId, text)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return msg
}
