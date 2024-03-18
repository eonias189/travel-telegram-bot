package tgrouter

import (
	"fmt"

	errlogger "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/loggers/errLogger"
)

type Config struct {
	OnNotFound HandlerFunc
	OnText     HandlerFunc
	OnError    func(error)
	Workers    int
}

var DefaultConfig = &Config{
	OnNotFound: func(ctx *Context) error {
		return ctx.SendString(fmt.Sprintf("command not found: %v", ctx.Update.Message.Command()))
	},
	OnError: func(err error) {
		errlogger.New().Error(err.Error())
	},
	OnText: func(ctx *Context) error {
		return nil
	},
	Workers: 1,
}

func handleConfig(cfg *Config) *Config {
	if cfg == nil {
		return DefaultConfig
	}
	if cfg.OnError == nil {
		cfg.OnError = DefaultConfig.OnError
	}
	if cfg.OnNotFound == nil {
		cfg.OnNotFound = DefaultConfig.OnNotFound
	}
	if cfg.OnText == nil {
		cfg.OnText = DefaultConfig.OnText
	}
	if cfg.Workers == 0 {
		cfg.Workers = DefaultConfig.Workers
	}
	return cfg
}
