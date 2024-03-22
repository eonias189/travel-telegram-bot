package msgtempl

import (
	"fmt"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ProfileMessage(name string, user service.User) string {
	return fmt.Sprintf(`Профиль:
Пользователь: %v
Возраст: %v
Страна: %v
Город: %v
bio: %v`, name, user.Age, user.Country, user.City, user.Bio)
}

func ProfileButtons() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("назад", "menu")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить возраст", "change-age")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить город", "change-city")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить страну", "change-country")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("изменить bio", "change-bio")),
	)
}
