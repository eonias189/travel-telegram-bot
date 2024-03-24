package msgtempl

import (
	"fmt"
	"time"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/geoapi"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LocationsMessage(senderId, tripId int64, locations []service.Location) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(senderId, "локации:")
	rows := make([][]tgbotapi.InlineKeyboardButton, len(locations)+2)

	for i, location := range locations {
		rows[i+2] = tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(`%v. %v`, location.Order, location.Name),
				fmt.Sprintf(`location?tripId=%v&order=%v`, tripId, location.Order),
			),
		)
	}
	rows[0] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", fmt.Sprintf(`trip?id=%v`, tripId)))
	rows[1] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("добавить", fmt.Sprintf(`new-location?tripId=%v`, tripId)))

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(rows...)
	return msg
}

func LocationMessage(tripId, senderId int64, location service.Location) tgbotapi.MessageConfig {
	coords, err := geoapi.GetCoords(location.Name)
	var coordsStr string

	if err != nil {
		coordsStr = "не найдено"
	} else {
		coordsStr = fmt.Sprintf(`%v, %v`, coords.Lat, coords.Lng)
	}

	startTime := time.Unix(0, location.StartTime)
	endTime := time.Unix(0, location.EndTime)

	text := fmt.Sprintf(`%v
Координаты: %v
Дата прибытия: %v
Дата отбытия: %v`, location.Name, coordsStr, startTime.Format(`02.01.2006/15:04`), endTime.Format(`02.01.2006/15:04`))

	msg := tgbotapi.NewMessage(senderId, text)

	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", fmt.Sprintf(`locations?tripId=%v`, tripId))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить локацию", fmt.Sprintf(`change-location-name?tripId=%v&order=%v`, tripId, location.Order))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить дату прибытя", fmt.Sprintf(`change-location-start?tripId=%v&order=%v`, tripId, location.Order))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить дату отбытия", fmt.Sprintf(`change-location-end?tripId=%v&order=%v`, tripId, location.Order))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("удалить", fmt.Sprintf(`delete-location?tripId=%v&order=%v`, tripId, location.Order))),
	)

	msg.ReplyMarkup = btns
	return msg
}
