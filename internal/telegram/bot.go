package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/internal/game"
	"time"
)

type Bot struct {
	bot   *tgbotapi.BotAPI
	game  *game.Game
	speed time.Duration
}

func NewBot(bot *tgbotapi.BotAPI, speed time.Duration) *Bot {
	return &Bot{
		bot:   bot,
		speed: speed,
	}
}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			if err := b.handleCallback(update.CallbackQuery); err != nil {
				fmt.Println(err.Error())
			}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				fmt.Println(err.Error())
			}

			continue
		}

		// Handle regular messages
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if _, err = b.bot.Send(msg); err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.StopReceivingUpdates()

	return nil
}
