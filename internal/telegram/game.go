package telegram

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
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
	if !b.games[gameId].Finished() {
		return nil
	}

	b.drawLast(gameId)

	ctx := context.Background()
	toCreate := domain.GameToCreate{
		ID:     gameId.String(),
		Score1: b.games[gameId].Score.First,
		Score2: b.games[gameId].Score.Second,
		Rounds: b.games[gameId].RoundsCount,
	}

	players := make([]domain.PlayerTeam, 0)

	for i, p := range b.games[gameId].Players {
		if p.Bot {
			continue
		}

		players = append(players, domain.PlayerTeam{
			ID:   p.ID,
			Team: uint8(i%2 + 1),
		})
	}

	if _, err = b.services.Games.Save(ctx, toCreate, players...); err != nil {
		return err
	}

	b.log.Info("game with id " + gameId.String() + " saved")

	return squirrel.ErrGameFinished
}

func (b *Bot) draw(gameId uuid.UUID) {
	ui := b.drawGame(gameId, 0, false)
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

func (b *Bot) drawLast(gameId uuid.UUID) {
	ui := b.drawGame(gameId, 0, true)
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
