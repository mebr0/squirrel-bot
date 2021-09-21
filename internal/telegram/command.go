package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case "start":
		return b.handleStartCommand(message)
	case "games":
		return b.handleGamesCommand(message)
	case "game":
		return b.handleGameCommand(message)
	}

	return nil
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	p := domain.NewPlayer(message.Chat.ID, message.Chat.UserName, message.Chat.FirstName, message.Chat.LastName)

	_, err := b.services.Players.Register(context.Background(), p)

	if err != nil {
		return err
	}

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: message.Chat.ID,
		},
		Text: "Hello",
	}

	_, err = b.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleGamesCommand(message *tgbotapi.Message) error {
	ctx := context.Background()

	games, err := b.services.Games.List(ctx, message.Chat.ID)

	if err != nil {
		return err
	}

	ui := b.lineGames(games)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID: message.Chat.ID,
		},
		Text: ui,
	}

	_, err = b.bot.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func (b *Bot) handleGameCommand(message *tgbotapi.Message) error {
	player := squirrel.NewPlayer(message.Chat.ID, message.Chat.UserName, message.Chat.FirstName, message.Chat.LastName)

	id := uuid.New()

	b.games[id] = squirrel.NewGameWithBots(player)

	b.games[id].StartFirstRound()

	ui := b.drawGame(id, 0, false)

	msg := tgbotapi.MessageConfig{
		BaseChat: tgbotapi.BaseChat{
			ChatID:      message.Chat.ID,
			ReplyMarkup: b.inlineKeyboard(id, 0),
		},
		Text:                  ui,
		ParseMode:             "MarkdownV2",
		DisableWebPagePreview: false,
	}

	m, err := b.bot.Send(msg)

	if err != nil {
		return err
	}

	b.games[id].UpdateChats([4]int{m.MessageID, -1, -1, -1})

	return nil
}
