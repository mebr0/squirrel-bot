package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"strconv"
	"strings"
)

const (
	// Common constants for callback data
	callbackPrefix        = "sqrl"
	callbackTextSeparator = "_"

	// Game constants for callback data
	gameCallback            = "game"
	gameCallbackThrowAction = "throw"
)

func (b *Bot) handleCallback(callback *tgbotapi.CallbackQuery) error {
	tokens := strings.Split(callback.Data, callbackTextSeparator)

	// If callback data does not contain general prefix, ignore it
	if len(tokens) < 1 || tokens[0] != callbackPrefix {
		return ErrUnknownCallback
	}

	switch tokens[1] {
	case gameCallback:
		return b.handleGamesCallback(callback)
	default:
		return ErrUnknownCallback
	}
}

func (b *Bot) handleGamesCallback(callback *tgbotapi.CallbackQuery) error {
	tokens := strings.Split(callback.Data, callbackTextSeparator)

	// If callback data does not contain game id, ignore it
	if len(tokens) < 2 {
		return ErrUnknownCallback
	}

	gameId, err := uuid.Parse(tokens[2])

	if err != nil {
		return nil
	}

	return b.handleGameCallback(callback, gameId)
}

func (b *Bot) handleGameCallback(callback *tgbotapi.CallbackQuery, gameId uuid.UUID) error {
	tokens := strings.Split(callback.Data, callbackTextSeparator)

	// If callback data does not contain game action, ignore it
	if len(tokens) < 3 {
		return ErrUnknownCallback
	}

	switch tokens[3] {
	case gameCallbackThrowAction:
		return b.handleGameCallbackThrowActions(callback, gameId)
	default:
		return ErrUnknownCallback
	}
}

func (b *Bot) handleGameCallbackThrowActions(callback *tgbotapi.CallbackQuery, gameId uuid.UUID) error {
	tokens := strings.Split(callback.Data, callbackTextSeparator)

	// If callback data does not contain card to throw, ignore it
	if len(tokens) < 4 {
		return ErrUnknownCallback
	}

	card, err := strconv.ParseUint(tokens[4], 10, 8)

	if err != nil {
		return err
	}

	cardN := squirrel.Card(uint8(card))

	if err = b.games[gameId].Throw(cardN); err != nil {
		return err
	}

	b.draw(gameId)

	if err = b.processGame(gameId); err != nil {
		if err == squirrel.ErrGameFinished {
			return nil
		}

		return err
	}

	for b.games[gameId].BotsTurn() {
		if err = b.games[gameId].BotMove(); err != nil {
			return err
		}

		b.draw(gameId)

		if err = b.processGame(gameId); err != nil {
			if err == squirrel.ErrGameFinished {
				return nil
			}

			return err
		}
	}

	return nil
}

func throwCardCallback(gameId uuid.UUID, card squirrel.Card) string {
	// Convert card to string representation, because it will not work in newCallback
	cardString := strconv.FormatInt(int64(card), 10)

	return newCallback(callbackPrefix, gameCallback, gameId, gameCallbackThrowAction, cardString)
}

func newCallback(tokens ...interface{}) string {
	values := make([]string, len(tokens))

	for i, v := range tokens {
		values[i] = fmt.Sprintf("%s", v)
	}

	return strings.Join(values, callbackTextSeparator)
}
