package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Central-University-IT-prod/backend-eonias189/internal/app"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/config"
	applogger "github.com/Central-University-IT-prod/backend-eonias189/internal/lib/loggers/appLogger"
	"github.com/Central-University-IT-prod/backend-eonias189/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}

	logger := applogger.New(cfg.Env)

	rdb, err := service.Connect(ctx, cfg.RedisAddr, cfg.RedisUser, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(0)
	}

	app := app.New(rdb, logger)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		cancel()
		if rdb.Save(context.TODO()).Err() == nil {
			logger.Info("data saved")
		}
	}()

	logger.Info("connected to redis")

	err = app.Run(ctx, cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
}
