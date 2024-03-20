package config

import (
	"fmt"
	"os"
	"strconv"
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
	Env           LoggerEnv
	BotToken      string
	RedisAddr     string
	RedisUser     string
	RedisPassword string
	RedisDB       int
}

func Get() (*Config, error) {
	env := os.Getenv("ENV")

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, ErrMissingEnvVar("BOT_TOKEN")
	}

	rhaddr := os.Getenv("REDIS_ADDRESS")
	if rhaddr == "" {
		return nil, ErrMissingEnvVar("REDIS_ADDRESS")
	}

	ruser := os.Getenv("REDIS_USER")

	rpassword := os.Getenv("REDIS_PASSWORD")

	rdb := os.Getenv("REDIS_DB")
	rdbInt, err := strconv.Atoi(rdb)
	if err != nil {
		return nil, err
	}

	return &Config{
		Env:           LoggerEnv(env),
		BotToken:      token,
		RedisAddr:     rhaddr,
		RedisUser:     ruser,
		RedisPassword: rpassword,
		RedisDB:       rdbInt,
	}, nil
}
