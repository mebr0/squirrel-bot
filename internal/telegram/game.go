package telegram

import (
	"fmt"
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
	ui := fmt.Sprintf(`Score: %d:%d | Round Score: %d:%d
Turn: %d | Trump: %d
Board: %s %s %s %s
P2: %s
P3: %s
P4: %s`, b.game.Score.First, b.game.Score.Second, b.game.Board.Round.First, b.game.Board.Round.Second,
		b.game.Board.CurrentTurn, b.game.Board.Trump, b.game.Board.Cards[0].Symbol(), b.game.Board.Cards[1].Symbol(),
		b.game.Board.Cards[2].Symbol(), b.game.Board.Cards[3].Symbol(), b.game.Players[1].Hand.ToString(), b.game.Players[2].Hand.ToString(), b.game.Players[3].Hand.ToString())

	msg := tgbotapi.NewEditMessageText(b.game.Players[0].ID, b.game.Players[0].Message, ui)
	_, err := b.bot.Send(msg)

	if err != nil {
		b.log.Error("error editing message - " + err.Error())
	}

	mark := tgbotapi.NewEditMessageReplyMarkup(b.game.Players[0].ID, b.game.Players[0].Message,
		inlineKeyboard(*b.game.Players[0].Hand))

	_, err = b.bot.Send(mark)

	if err != nil {
		b.log.Error("error editing message inline keyboard" + err.Error())
	}

	time.Sleep(b.speed)
}
