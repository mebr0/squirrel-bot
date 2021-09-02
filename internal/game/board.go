package game

type Board struct {
	Cards [playersCount]Card
	Round *RoundScore

	Trump CardSuit

	FirstTurn   uint8
	CurrentTurn uint8

	BottomIndex uint8
}

const size uint8 = 4

func newBoard() *Board {
	return &Board{
		Cards: [playersCount]Card{
			emptyCard,
			emptyCard,
			emptyCard,
			emptyCard,
		},
		Round:       newRoundScore(),
		BottomIndex: size,
		Trump:       clubs,
		FirstTurn:   0,
		CurrentTurn: 0,
	}
}

func (b *Board) refresh(roundCount uint8) {
	for i := range b.Cards {
		b.Cards[i] = emptyCard
	}

	b.BottomIndex = size

	b.CurrentTurn = roundCount % size
}

func (b *Board) bottom() (Card, error) {
	if b.BottomIndex == size {
		return b.Cards[0], ErrEmptyBoard
	}

	return b.Cards[b.BottomIndex], nil
}

func (b *Board) BottomSuit() (CardSuit, error) {
	bottom, err := b.bottom()

	if err != nil {
		return 0, err
	}

	if bottom.IsEmpty() {
		return 0, ErrEmptyBoard
	}

	return b.Cards[b.BottomIndex].suit(), nil
}

func (b *Board) throw(index uint8, card Card) {
	b.Cards[index] = card

	if b.BottomIndex == size {
		b.BottomIndex = index
	}

	b.nextTurn()
}

func (b *Board) nextTurn() {
	b.CurrentTurn = (b.CurrentTurn + 1) % size
}

func (b *Board) isFull() bool {
	for _, c := range b.Cards {
		if c == emptyCard {
			return false
		}
	}

	return true
}

func (b *Board) index(card Card) uint8 {
	for i, c := range b.Cards {
		if c == card {
			return uint8(i)
		}
	}

	return 4
}

func (b *Board) sum() uint8 {
	var sum uint8

	for _, c := range b.Cards {
		sum += c.Points()
	}

	return sum
}

func (b *Board) empty() {
	for i := range b.Cards {
		b.Cards[i] = emptyCard
	}

	b.BottomIndex = size
}

func (b *Board) calculate() error {
	greatest := b.greatestCard()
	greatestIndex := b.index(greatest)

	if err := b.Round.add(b.sum(), Team(greatestIndex%2+1)); err != nil {
		return err
	}

	b.empty()

	b.CurrentTurn = greatestIndex

	return nil
}

func (b *Board) greatestCard() Card {
	bottom, err := b.bottom()

	if err != nil {
		return emptyCard
	}

	greatest := emptyCard

	if bottom.isTrump(b.Trump) {
		for _, c := range b.Cards {
			if greatest == emptyCard {
				greatest = c
				continue
			}

			if c.value() == jacks {
				if greatest.value() == jacks && c > greatest {
					greatest = c
				}

				if greatest.value() != jacks {
					greatest = c
				}

				continue
			}

			if c.value() != jacks {
				if greatest.value() != jacks && c > greatest {
					greatest = c
				}
			}
		}
	}

	if !bottom.isTrump(b.Trump) {
		suit := bottom.suit()

		for _, c := range b.Cards {
			if c.suit() == suit && !c.isTrump(b.Trump) {
				if greatest == emptyCard {
					greatest = c
					continue
				}

				if !greatest.isTrump(b.Trump) && c > greatest {
					greatest = c
					continue
				}
			}

			if c.isTrump(b.Trump) {
				if greatest == emptyCard {
					greatest = c
					continue
				}

				if !greatest.isTrump(b.Trump) {
					greatest = c
					continue
				}

				if greatest.isTrump(b.Trump) {
					if c.value() == jacks {
						if greatest.value() == jacks && c > greatest {
							greatest = c
						}

						if greatest.value() != jacks {
							greatest = c
						}

						continue
					}

					if c.value() != jacks {
						if greatest.value() != jacks {
							greatest = c
							continue
						}
					}
				}
			}
		}
	}

	return greatest
}

func (b *Board) roundFinished() bool {
	return b.Round.finished()
}
