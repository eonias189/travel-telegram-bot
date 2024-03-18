package texthandler

import "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/tgrouter"

type Config struct {
	OnUnknownContext tgrouter.HandlerFunc
}

var defaultConfig = &Config{
	OnUnknownContext: func(ctx *tgrouter.Context) error {
		return nil
	},
}

func handleConfig(cfg *Config) *Config {
	if cfg == nil {
		return defaultConfig
	}

	if cfg.OnUnknownContext == nil {
		cfg.OnUnknownContext = defaultConfig.OnUnknownContext
	}

	return cfg
}
