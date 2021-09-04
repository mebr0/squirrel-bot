package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"strconv"
	"strings"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	tokens := strings.Split(callback.Data, "_")
	card, err := strconv.ParseUint(tokens[2], 10, 8)

	if err != nil {
		return err
	}

	cardN := squirrel.Card(uint8(card))

	if err = b.game.Throw(cardN); err != nil {
		return err
	}

	b.draw()

	if err = b.processGame(); err != nil {
		return err
	}

	for b.game.BotsTurn() {
		if err = b.game.BotMove(); err != nil {
			return err
		}

		b.draw()

		if err = b.processGame(); err != nil {
			return err
		}
	}

	return nil
}
