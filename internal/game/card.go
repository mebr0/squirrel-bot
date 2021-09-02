package game

import (
	"math/rand"
	"time"
)

// CardValue represents values between 7 and 10 plus queens, kings, aces and jacks
type CardValue uint8

const (
	sevens CardValue = iota
	eights
	nines
	queens
	kings
	tens
	aces
	jacks
)

// CardSuit represents suits: ♦ ♥ ♠ ♣. Ordered by increasing jacks' power
type CardSuit uint8

const (
	diamonds CardSuit = iota
	hearts
	spades
	clubs
)

// Card represents any card deck of 32 cards which are combinations of all CardValue and CardSuit
type Card uint8

const (
	emptyCard Card = 32
)

func (c Card) IsEmpty() bool {
	return c == emptyCard
}

func (c Card) number() uint8 {
	return uint8(c)
}

func (c Card) value() CardValue {
	return CardValue(c.number() % handSize)
}

func (c Card) suit() CardSuit {
	return CardSuit(c.number() / handSize)
}

// Card is trump, if its suit matches with trump, or it is jacks
func (c Card) isTrump(trump CardSuit) bool {
	return c.suit() == trump || c.value() == jacks
}

// Every Card have own price represented by points
var points = map[CardValue]uint8{
	sevens: 0,
	eights: 0,
	nines:  0,
	jacks:  2,
	queens: 3,
	kings:  4,
	tens:   10,
	aces:   11,
}

func (c Card) Points() uint8 {
	return points[c.value()]
}

// Init fresh deck with ordered cards
func deck() []Card {
	cards := make([]Card, emptyCard)

	for i := range cards {
		cards[i] = Card(i)
	}

	return cards
}

// Shuffle and return new deck
func shuffledDeck() []Card {
	cards := deck()

	rand.Seed(time.Now().UnixNano())

	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})

	return cards
}

func (c Card) Symbol() string {
	if c == emptyCard {
		return "EM"
	}

	result := ""

	switch c.value() {
	case sevens:
		result += "7"
	case eights:
		result += "8"
	case nines:
		result += "9"
	case tens:
		result += "10"
	case jacks:
		result += "J"
	case queens:
		result += "Q"
	case kings:
		result += "K"
	case aces:
		result += "A"
	}

	switch c.suit() {
	case spades:
		result += "♠"
	case clubs:
		result += "♣"
	case hearts:
		result += "♥"
	case diamonds:
		result += "♦"
	}

	return result
}
