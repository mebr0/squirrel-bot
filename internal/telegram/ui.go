package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"strconv"
	"strings"
)

const (
	width = 40
	boardUnit = 7
)

func (b *Bot) drawGame(playerIndex int) string {
	team := squirrel.Team(playerIndex%2 + 1)

	score := b.game.Score.String(team)
	trump := b.game.Board.Trump.String()
	roundsCount := b.game.RoundsCount
	roundScore := b.game.Board.Round.String(team)
	players := b.game.Players.Shifted(playerIndex)
	cards := b.game.Board.ShiftedCards(playerIndex)

	ui := "```"
	row := ""

	row += fmt.Sprintf("Score: %s | Trump: %s | Round: %d", score, trump, roundsCount)
	row += strings.Repeat(" ", width - len(row)) + "\n"
	ui += row

	row = ""
	row += fmt.Sprintf("Round Score: %s", roundScore)
	row += strings.Repeat(" ", width - len(row)) + "\n"

	ui += row + "\n"
	ui += strings.Repeat("-", width) + "\n"

	ui += alignCenter(players[2].NickName()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(cards[2].Symbol()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"

	ui += spaceBetween(players[1].NickName(), cards[1].Symbol(), cards[3].Symbol(), players[3].NickName()) + "\n"

	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(cards[0].Symbol()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(players[0].NickName()) + "\n"

	ui += "```"

	return ui
}

func alignCenter(text string) string {
	spaces := strings.Repeat(" ", (width - len(text)) / 2)

	return spaces + text + spaces
}

func spaceBetween(nickName1, card1, card3, nickName3 string) string {
	result := nickName1 + strings.Repeat(" ", width / 2 - len(nickName1) - len([]rune(card1)) - boardUnit) + card1
	result += strings.Repeat(" ", 2 * boardUnit)
	result += card3 + strings.Repeat(" ", width / 2 - len(nickName3) - len([]rune(card3)) - boardUnit) + nickName3

	return result
}

func (b *Bot) inlineKeyboard(index int) tgbotapi.InlineKeyboardMarkup {
	cards := b.game.Players[index].Hand

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