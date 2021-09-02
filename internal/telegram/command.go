package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/internal/game"
	"strconv"
)

func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case "game":
		return b.handleStartCommand(message)
	}

	return nil
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	b.game = game.NewGame([4]int64{message.Chat.ID, -1, -1, -1})

	b.game.StartFirstRound()

	ui := fmt.Sprintf(`Score: %d:%d | Round Score: %d:%d
Turn: %d | Trump: %d
Board: %s %s %s %s
P2: %s
P3: %s
P4: %s`, b.game.Score.First, b.game.Score.Second, b.game.Board.Round.First, b.game.Board.Round.Second,
		b.game.Board.CurrentTurn, b.game.Board.Trump, b.game.Board.Cards[0].Symbol(), b.game.Board.Cards[1].Symbol(),
		b.game.Board.Cards[2].Symbol(), b.game.Board.Cards[3].Symbol(),
		b.game.Players[1].Hand.ToString(), b.game.Players[2].Hand.ToString(), b.game.Players[3].Hand.ToString())

	msg := tgbotapi.NewMessage(message.Chat.ID, ui)
	m, err := b.bot.Send(msg)

	if err != nil {
		return err
	}

	b.game.UpdateChats([4]int{m.MessageID, -1, -1, -1})

	mark := tgbotapi.NewEditMessageReplyMarkup(message.Chat.ID, m.MessageID, inlineKeyboard(*b.game.Players[0].Hand))

	if _, err := b.bot.Send(mark); err != nil {
		return err
	}

	return nil
}

func inlineKeyboard(cards game.Hand) tgbotapi.InlineKeyboardMarkup {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i := 0; i < 2; i++ {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{})

		for j := 0; j < 4; j++ {
			card := cards[i*4+j]

			if !card.IsEmpty() {
				callback := "sqrl_throw_" + strconv.FormatUint(uint64(card), 10)

				keyboard[i] = append(keyboard[i], tgbotapi.InlineKeyboardButton{
					Text:         card.Symbol(),
					CallbackData: &callback,
				})
			}
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
