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
)

func main() {
	cfg, err := config.Get()

	if err != nil {
		log.Fatal(err)
	}

	logger := applogger.New(cfg.Env)
	app := app.New(logger)
	ctx, cancel := context.WithCancel(context.Background())

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sigCh
		cancel()
		app.Close()
	}()

	err = app.Run(ctx, cfg.BotToken)
	if err != nil {
		log.Fatal(err)
	}
}
