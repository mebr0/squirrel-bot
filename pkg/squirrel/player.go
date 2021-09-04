package squirrel

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

func (h Hand) ToString() string {
	res := ""

	for _, c := range h {
		res += c.Symbol() + " "
	}

	return res
}

type Players [playersCount]*Player

type Player struct {
	ID      int64
	Message int
	Hand    *Hand
	Bot     bool
}

const (
	playersCount uint8 = 4
)

func newPlayer(id int64) *Player {
	return &Player{
		ID:   id,
		Hand: emptyHand(),
		Bot:  false,
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
