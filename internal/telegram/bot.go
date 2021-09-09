package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mebr0/squirrel-bot/internal/config"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"go.uber.org/zap"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	games  map[uuid.UUID]*squirrel.Game
	log    *zap.Logger
	config config.Game
}

func NewBot(bot *tgbotapi.BotAPI, log *zap.Logger, config config.Game) *Bot {
	return &Bot{
		bot:    bot,
		games:  map[uuid.UUID]*squirrel.Game{},
		log:    log,
		config: config,
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
				b.log.Error("error handling callback query - " + err.Error())
			}
		}

		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		// Handle commands
		if update.Message.IsCommand() {
			if err := b.handleCommand(update.Message); err != nil {
				b.log.Error("error handling command - " + err.Error())
			}

			continue
		}

		// Handle regular messages
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if _, err = b.bot.Send(msg); err != nil {
			b.log.Error("error handling message text - " + err.Error())
		}
	}

	return nil
}

func (b *Bot) Stop(ctx context.Context) error {
	b.bot.StopReceivingUpdates()

	return nil
}
