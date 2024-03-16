package config

import (
	"fmt"
	"os"
)

type LoggerEnv string

var (
	LocalEnv       LoggerEnv = "local"
	DevelopmentEnv LoggerEnv = "development"
	ProductionEnv  LoggerEnv = "production"
)

func ErrMissingEnvVar(variable string) error {
	return fmt.Errorf("MISSING REQUIRED ENV VARIABLE: %v", variable)
}

type Config struct {
	Env      LoggerEnv
	BotToken string
}

func Get() (*Config, error) {
	cfg := &Config{}

	env := os.Getenv("ENV")
	cfg.Env = LoggerEnv(env)

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, ErrMissingEnvVar("BOT_TOKEN")
	}

	cfg.BotToken = token
	return cfg, nil
}
