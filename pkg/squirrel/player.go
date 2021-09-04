package squirrel

import "strings"

// Hand is collection of Card of one single Player
type Hand [handSize]Card

const (
	handSize uint8 = 8
)

func emptyHand() *Hand {
	return &Hand{
		emptyCard,
		emptyCard,
		emptyCard,
		emptyCard,
		emptyCard,
		emptyCard,
		emptyCard,
		emptyCard,
	}
}

type Players [playersCount]*Player

type Player struct {
	ID       int64
	Message  int
	Username string
	Name     string
	LastName string
	Hand     *Hand
	Bot      bool
}

const (
	playersCount uint8 = 4
)

func newWithBots(player *Player) *Players {
	var players [playersCount]*Player

	players[0] = player
	players[1] = newBot()
	players[2] = newBot()
	players[3] = newBot()

	return (*Players)(&players)
}

func (p *Players) Shifted(index int) *Players {
	shifted := append(p[index:], p[:index]...)

	var players [playersCount]*Player

	copy(players[:], shifted)

	return (*Players)(&players)
}

func NewPlayer(id int64, username, name, lastName string) *Player {
	return &Player{
		ID:       id,
		Username: username,
		Name:     name,
		LastName: lastName,
		Hand:     emptyHand(),
		Bot:      false,
	}
}

func newBot() *Player {
	return &Player{
		Hand: emptyHand(),
		Bot:  true,
	}
}

func (p *Players) playerIndexByCard(card Card) uint8 {
	for i, player := range p {
		for _, c := range player.Hand {
			if c == card {
				return uint8(i)
			}
		}
	}

	return playersCount
}

func (p *Players) dealCards() {
	cards := shuffledDeck()

	for i := range p {
		var hand Hand

		copy(hand[:], cards[i*8:(i+1)*8])

		p[i].Hand = &hand
	}
}

// Check whether player can play anything that matches bottom card
func (p *Player) matchBottomCard(bottomCard Card, trump CardSuit) bool {
	if bottomCard.isTrump(trump) {
		for _, c := range p.Hand {
			if c.isTrump(trump) {
				return true
			}
		}

		return false
	}

	suit := bottomCard.suit()

	for _, c := range p.Hand {
		if c.suit() == suit && !c.isTrump(trump) {
			return true
		}
	}

	return false
}

func (p *Player) throw(card Card) {
	for i, c := range p.Hand {
		if c == card {
			p.Hand[i] = emptyCard
			break
		}
	}
}

func (p *Player) NickName() string {
	if p.Bot {
		return "Bot"
	}

	if p.Username != "" {
		return "@" + p.Username
	}

	fullName := p.Name + " " + p.LastName

	return strings.TrimSpace(fullName)
}
