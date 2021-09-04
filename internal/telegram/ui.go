package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mebr0/squirrel-bot/pkg/squirrel"
	"strconv"
)

func (b *Bot) drawGame(playerIndex int) string {
	team := squirrel.Team(playerIndex%2 + 1)

	score := b.game.Score.String(team)
	trump := b.game.Board.Trump.String()
	roundsCount := b.game.RoundsCount
	roundScore := b.game.Board.Round.String(team)
	players := b.game.Players.Shifted(playerIndex)
	cards := b.game.Board.ShiftedCards(playerIndex)

	return fmt.Sprintf(`
Score: %s | Trump: %s | Round: %d
Round Score: %s
—————————————————
                           %s
                                      
                              %s

    %s         %s               %s️   %s    

                              %s
                     
                             %s`, score, trump, roundsCount, roundScore, players[2].NickName(), cards[2].Symbol(),
		players[1].NickName(), cards[1].Symbol(), cards[3].Symbol(), players[3].NickName(),
		cards[0].Symbol(), players[0].NickName())
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
