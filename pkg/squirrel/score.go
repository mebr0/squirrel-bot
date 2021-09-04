package squirrel

type Team uint8

const (
	noTeam Team = iota
	firstTeam
	secondTeam
	draw
)

type BaseScore interface {
	add(value uint8, team uint8) error
	finished() bool
	winner() Team
	refresh()
}

const (
	minPerTime      = 1
	maxPerTime      = 4
	overallPoints   = 23
	upperBoundToWin = 11
	maxScore        = 12
)

type Score struct {
	First  uint8
	Second uint8
}

func newScore() *Score {
	return &Score{
		First:  0,
		Second: 0,
	}
}

func (s *Score) add(value uint8, team Team) error {
	if value < minPerTime || value > maxPerTime {
		return ErrInvalidScore
	}

	if !s.checkAddValue(value) {
		return ErrScoreExceeded
	}

	if team == firstTeam {
		score := s.First + value

		if score > maxScore {
			score = maxScore
		}

		s.First = score

		return nil
	}

	if team == secondTeam {
		score := s.Second + value

		if score > maxScore {
			score = maxScore
		}

		s.Second = score

		return nil
	}

	return ErrInvalidTeam
}

func (s *Score) checkAddValue(value uint8) bool {
	return s.First+value+s.Second <= overallPoints
}

func (s *Score) finished() bool {
	return s.First > upperBoundToWin || s.Second > upperBoundToWin
}

//func (s *Score) winner() Team {
//	if s.First > upperBoundToWin {
//		return firstTeam
//	}
//
//	if s.Second > upperBoundToWin {
//		return secondTeam
//	}
//
//	return noTeam
//}

//func (s *Score) refresh() {
//	s.First = 0
//	s.Second = 0
//}

const (
	roundMinPerTime      = 0
	roundMaxPerTime      = 44
	roundOverallPoints   = 120
	roundUpperBoundToWin = 60
)

type RoundScore struct {
	First  uint8
	Second uint8
}

func newRoundScore() *RoundScore {
	return &RoundScore{
		First:  0,
		Second: 0,
	}
}

func (s *RoundScore) add(value uint8, team Team) error {
	if value < roundMinPerTime || value > roundMaxPerTime {
		return ErrInvalidScore
	}

	if !s.checkAddValue(value) {
		return ErrScoreExceeded
	}

	if team == firstTeam {
		s.First += value
		return nil
	}

	if team == secondTeam {
		s.Second += value
		return nil
	}

	return ErrInvalidTeam
}

func (s *RoundScore) checkAddValue(value uint8) bool {
	return s.First+value+s.Second <= roundOverallPoints
}

func (s *RoundScore) finished() bool {
	return s.First+s.Second == roundOverallPoints
}

func (s *RoundScore) winner() Team {
	if s.First == s.Second && s.First*2 == roundOverallPoints {
		return draw
	}

	if s.First > roundUpperBoundToWin {
		return firstTeam
	}

	if s.Second > roundUpperBoundToWin {
		return secondTeam
	}

	return noTeam
}

func (s *RoundScore) refresh() {
	s.First = 0
	s.Second = 0
}
