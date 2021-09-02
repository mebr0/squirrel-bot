package app

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/internal/config"
	"github.com/mebr0/squirrel-bot/internal/telegram"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(configPath string) {
	cfg := config.LoadConfig(configPath)

	botApi, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)

	if err != nil {
		log.Fatal(err)
	}

	//botApi.Debug = true

	// Telegram bot
	bot := telegram.NewBot(botApi, cfg.Game.Speed)
	go func() {
		if err = bot.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err = bot.Stop(ctx); err != nil {
		log.Println("failed to stop server: " + err.Error())
	}
}
