package squirrel

import "errors"

var (
	ErrInvalidScore  = errors.New("invalid score")
	ErrScoreExceeded = errors.New("invalid score, overall score exceeded")
	ErrInvalidTeam   = errors.New("invalid team to add score")

	ErrGameFinished    = errors.New("game already finished")
	ErrGameNotFinished = errors.New("game not finished")

	ErrEmptyBoard = errors.New("empty board")

	ErrCardNotFound = errors.New("card not found")
	ErrNotYourTurn  = errors.New("not your turn")

	ErrMustThrowTrump     = errors.New("must throw trump")
	ErrMustThrowOtherCard = errors.New("must throw other card")

	ErrNoValidMoves = errors.New("no valid moves")
)
