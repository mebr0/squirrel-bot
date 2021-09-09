package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"time"
)

func (b *Bot) processGame(gameId uuid.UUID) error {
	// Check board full
	fullBoard, err := b.games[gameId].BoardFull()

	if err != nil {
		return err
	}

	if !fullBoard {
		return nil
	}

	b.draw(gameId)

	// Check round finished
	roundFinished, err := b.games[gameId].RoundFinished()

	if err != nil {
		return err
	}

	if !roundFinished {
		return nil
	}

	b.draw(gameId)

	// Check game finished
	gameFinished, err := b.games[gameId].Finished()

	if err != nil {
		return err
	}

	if !gameFinished {
		return nil
	}

	b.draw(gameId)

	return nil
}

func (b *Bot) draw(gameId uuid.UUID) {
	ui := b.drawGame(gameId, 0)
	keyboard := b.inlineKeyboard(gameId, 0)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.games[gameId].Players[0].ID,
			MessageID:   b.games[gameId].Players[0].Message,
			ReplyMarkup: &keyboard,
		},
		Text:      ui,
		ParseMode: "MarkdownV2",
	}

	_, err := b.bot.Send(msg)

	if err != nil {
		b.log.Error("error editing message - " + err.Error())
	}

	time.Sleep(b.config.Speed)
}
