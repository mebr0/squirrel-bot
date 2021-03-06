package squirrel

import (
	"math/rand"
	"time"
)

type Game struct {
	Players *Players
	Score   *Score
	Board   *Board

	RoundsCount uint8
}

func NewGameWithBots(player *Player) *Game {
	return &Game{
		Players: newWithBots(player),
		Score:   newScore(),
		Board:   newBoard(),

		RoundsCount: 1,
	}
}

func (g *Game) BotsTurn() bool {
	return g.Players[g.Board.CurrentTurn].Bot
}

func (g *Game) BotMove() error {
	if g.Finished() {
		return ErrGameFinished
	}

	cards, err := g.possibleTurns(g.Board.CurrentTurn)

	if err != nil {
		return err
	}

	if (len(cards)) == 0 {
		return ErrNoValidMoves
	}

	rand.Seed(time.Now().Unix())

	card := cards[rand.Intn(len(cards))]

	if err = g.Throw(card); err != nil {
		if err != ErrCardNotFound && err != ErrNotYourTurn {
			return err
		}
	}

	return nil
}

func (g *Game) UpdateChats(messageIDs [playersCount]int) {
	for i := range g.Players {
		g.Players[i].Message = messageIDs[i]
	}
}

func (g *Game) StartFirstRound() {
	g.Players.dealCards()
}

func (g *Game) Throw(card Card) error {
	if g.Finished() {
		return ErrGameFinished
	}

	playerIndex := g.Players.playerIndexByCard(card)

	if playerIndex == playersCount {
		return ErrCardNotFound
	}

	if playerIndex != g.Board.CurrentTurn {
		return ErrNotYourTurn
	}

	player := g.Players[playerIndex]

	if err := g.canThrow(playerIndex, card); err != nil {
		return err
	}

	g.Board.throw(playerIndex, card)
	player.throw(card)

	return nil
}

func (g *Game) possibleTurns(index uint8) ([]Card, error) {
	var cards []Card

	for _, c := range g.Players[index].Hand {
		if c == emptyCard {
			continue
		}

		err := g.canThrow(index, c)

		if err != nil {
			if err != ErrMustThrowTrump && err != ErrMustThrowOtherCard {
				return nil, err
			}

			continue
		}

		cards = append(cards, c)
	}

	return cards, nil
}

func (g *Game) canThrow(index uint8, card Card) error {
	bottom, err := g.Board.bottom()

	if err != nil {
		if err == ErrEmptyBoard {
			return nil
		}

		return err
	}

	trump := g.Board.Trump

	if g.Players[index].matchBottomCard(bottom, trump) {
		if bottom.isTrump(trump) {
			if !card.isTrump(trump) {
				return ErrMustThrowTrump
			}
		}

		if !bottom.isTrump(trump) {
			if !(card.suit() == bottom.suit() && !card.isTrump(trump)) {
				return ErrMustThrowOtherCard
			}
		}
	}

	return nil
}

func (g *Game) WinnerTeam() (Team, error) {
	if g.Score.finished() {
		return g.Score.winner(), nil
	}

	return 0, ErrGameNotFinished
}

func (g *Game) Finished() bool {
	return g.Score.finished()
}

func (g *Game) RoundFinished() (bool, error) {
	if g.Board.roundFinished() {
		if !g.Finished() {
			if err := g.nextRound(); err != nil {
				return true, err
			}
		}

		return true, nil
	}

	return false, nil
}

func (g *Game) nextRound() error {
	if err := g.updateScore(); err != nil {
		return err
	}

	g.Board.Round.refresh()

	g.Players.dealCards()

	g.RoundsCount += 1

	g.Board.refresh(g.RoundsCount - 1)

	return nil
}

func (g *Game) updateScore() error {
	winner := g.Board.Round.winner()

	return g.Score.add(2, winner)
}

func (g *Game) BoardFull() (bool, error) {
	if g.Board.isFull() {
		if err := g.Board.calculate(); err != nil {
			return true, err
		}

		return true, nil
	}

	return false, nil
}
