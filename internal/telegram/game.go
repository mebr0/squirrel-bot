package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"time"
)

func (b *Bot) processGame() error {
	// Check board full
	fullBoard, err := b.game.BoardFull()

	if err != nil {
		return err
	}

	if !fullBoard {
		return nil
	}

	b.draw()

	// Check round finished
	roundFinished, err := b.game.RoundFinished()

	if err != nil {
		return err
	}

	if !roundFinished {
		return nil
	}

	b.draw()

	// Check game finished
	gameFinished, err := b.game.Finished()

	if err != nil {
		return err
	}

	if !gameFinished {
		return nil
	}

	b.draw()

	return nil
}

func (b *Bot) draw() {
	ui := b.drawGame(0)
	keyboard := b.inlineKeyboard(0)

	msg := tgbotapi.EditMessageTextConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:      b.game.Players[0].ID,
			MessageID:   b.game.Players[0].Message,
			ReplyMarkup: &keyboard,
		},
		Text: ui,
		ParseMode: "MarkdownV2",
	}

	_, err := b.bot.Send(msg)

	if err != nil {
		b.log.Error("error editing message - " + err.Error())
	}

	time.Sleep(b.config.Speed)
}
