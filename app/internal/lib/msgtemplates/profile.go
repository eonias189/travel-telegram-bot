package msgtempl

import (
	"fmt"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProfileMsg(senderID int64, name string, user service.User) tgbotapi.MessageConfig {
	var age string
	if user.Age == -1 {
		age = "не указано"
	} else {
		age = fmt.Sprint(user.Age)
	}

	text := fmt.Sprintf(`Профиль:
Пользователь: @%v
Возраст: %v
Местоположение: %v
о себе: %v`, name, age, user.Location, user.Bio)
	btns := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "menu")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить возраст", "change-age")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить местоположение", "change-location")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(`изменить "о себе"`, "change-bio")),
	)

	msg := tgbotapi.NewMessage(senderID, text)
	msg.ReplyMarkup = btns
	return msg
}
