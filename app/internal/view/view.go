package view

import (
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type View interface {
	AddImage(name string, image tgbotapi.ReplyKeyboardMarkup)
	GetImage(name string) (tgbotapi.ReplyKeyboardMarkup, bool)
}

type viewImpl struct {
	images map[string]tgbotapi.ReplyKeyboardMarkup
	mu     sync.RWMutex
}

func (v *viewImpl) AddImage(name string, image tgbotapi.ReplyKeyboardMarkup) {
	v.mu.RLock()
	defer v.mu.RUnlock()
	v.images[name] = image
}

func (v *viewImpl) GetImage(name string) (tgbotapi.ReplyKeyboardMarkup, bool) {
	img, ok := v.images[name]
	return img, ok
}

func New() View {
	return &viewImpl{
		images: make(map[string]tgbotapi.ReplyKeyboardMarkup),
	}
}
