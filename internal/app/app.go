package app

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/internal/config"
	"github.com/mebr0/squirrel-bot/internal/telegram"
	"github.com/mebr0/squirrel-bot/pkg/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg := config.LoadConfig(configPath)

	// Logging
	log, err := logging.NewLogger(cfg.Log.Level)

	if err != nil {
		fmt.Println("error initializing logger instance - " + err.Error())
		return
	}

	//nolint:errcheck
	//goland:noinspection GoUnhandledErrorResult
	defer log.Sync()

	botApi, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)

	if err != nil {
		log.Fatal("error initializing telegram client - " + err.Error())
	}

	//botApi.Debug = true

	// Telegram bot
	bot := telegram.NewBot(botApi, log, cfg.Game.Speed)
	go func() {
		if err = bot.Start(); err != nil {
			log.Fatal("error while bot polling - " + err.Error())
		}
	}()

	log.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = bot.Stop(ctx); err != nil {
		log.Error("failed to stop server - " + err.Error())
	}
}
