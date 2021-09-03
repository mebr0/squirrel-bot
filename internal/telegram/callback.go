package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/internal/game"
	"strconv"
	"strings"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	tokens := strings.Split(callback.Data, "_")
	card, err := strconv.ParseUint(tokens[2], 10, 8)

	if err != nil {
		return err
	}

	cardN := game.Card(uint8(card))

	if err := b.game.Throw(cardN); err != nil {
		return err
	}

	b.draw()

	if err := b.processGame(); err != nil {
		return err
	}

	for b.game.BotsTurn() {
		b.game.BotMove()

		b.draw()

		if err := b.processGame(); err != nil {
			return err
		}
	}

	return nil
}