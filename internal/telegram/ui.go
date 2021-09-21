package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"strings"
)

const (
	width     = 40
	boardUnit = 7
)

func (b *Bot) lineGames(games []domain.Game) string {
	res := ""

	for i, g := range games {
		res += fmt.Sprintf("%d. Score: %d : %d. Rounds count: %d\n", i+1, g.Score1, g.Score2, g.Rounds)
	}

	return res
}

func (b *Bot) drawGame(gameId uuid.UUID, playerIndex int, finished bool) string {
	game := b.games[gameId]
	team := squirrel.Team(playerIndex%2 + 1)

	score := game.Score.String(team)
	trump := game.Board.Trump.String()
	roundsCount := game.RoundsCount
	roundScore := game.Board.Round.String(team)
	players := game.Players.Shifted(playerIndex)
	cards := game.Board.ShiftedCards(playerIndex)

	ui := "```\n"
	row := ""

	if finished {
		winnerTeam, err := game.WinnerTeam()

		if err != nil {
			b.log.Warn("Game not finished, while it supposed to - " + err.Error())
		} else {
			switch winnerTeam {
			case squirrel.FirstTeam:
				ui += "GAME FINISHED\n"
				row += fmt.Sprintf("Winners: %s and %s", game.Players[0].NickName(), game.Players[2].NickName())
			case squirrel.SecondTeam:
				ui += "GAME FINISHED\n"
				row += fmt.Sprintf("Winners: %s and %s", game.Players[1].NickName(), game.Players[3].NickName())
			default:
				b.log.Warn("Game not finished, while it supposed to - " + err.Error())
			}

			row += strings.Repeat(" ", width-len(row)) + "\n"
			ui += row
			ui += strings.Repeat("-", width) + "\n"
		}
	}

	row = ""
	row += fmt.Sprintf("Score: %s | Trump: %s | Round: %d", score, trump, roundsCount)
	row += strings.Repeat(" ", width-len(row)) + "\n"
	ui += row

	row = ""
	row += fmt.Sprintf("Round Score: %s", roundScore)
	row += strings.Repeat(" ", width-len(row)) + "\n"

	ui += row + "\n"
	ui += strings.Repeat("-", width) + "\n"

	ui += alignCenter(players[2].NickName()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(cards[2].Symbol()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"

	ui += spaceBetween(players[1].NickName(), cards[1], cards[3], players[3].NickName()) + "\n"

	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(cards[0].Symbol()) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += strings.Repeat(" ", width) + "\n"
	ui += alignCenter(players[0].NickName()) + "\n"

	ui += "```"

	return ui
}

func alignCenter(text string) string {
	spaces := strings.Repeat(" ", (width-len(text))/2)

	return spaces + text + spaces
}

func spaceBetween(nickName1 string, card1 squirrel.Card, card3 squirrel.Card, nickName3 string) string {
	result := nickName1 + strings.Repeat(" ", width/2-len(nickName1)-len([]rune(card1.Symbol()))-boardUnit) +
		card1.Symbol()
	result += strings.Repeat(" ", boardUnit)

	var thirdCardDelta = 0
	var boardDelta = 0

	if !card3.IsEmpty() {
		if !card1.IsEmpty() {
			boardDelta = 17
			thirdCardDelta = 9
			thirdCardDelta -= len(nickName3) - 3
		} else {
			thirdCardDelta = 7
			thirdCardDelta -= len(nickName3) - 3
		}
	} else {
		if !card1.IsEmpty() {
			boardDelta = 31
		}
	}

	result += strings.Repeat(" ", boardUnit+boardDelta)
	result += card3.Symbol() +
		strings.Repeat(" ", width/2-len(nickName3)-len([]rune(card3.Symbol()))-boardUnit+thirdCardDelta) +
		nickName3

	return result
}

func (b *Bot) inlineKeyboard(gameId uuid.UUID, index int) tgbotapi.InlineKeyboardMarkup {
	cards := b.games[gameId].Players[index].Hand

	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0)

	for i := 0; i < 2; i++ {
		keyboard = append(keyboard, []tgbotapi.InlineKeyboardButton{})

		for j := 0; j < 4; j++ {
			card := cards[i*4+j]

			if !card.IsEmpty() {
				callback := throwCardCallback(gameId, card)

				keyboard[i] = append(keyboard[i], tgbotapi.InlineKeyboardButton{
					Text:         card.Symbol(),
					CallbackData: &callback,
				})
			}
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}
