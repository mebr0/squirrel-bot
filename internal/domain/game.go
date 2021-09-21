package domain

import "time"

type Game struct {
	ID         string    `db:"id"`
	Score1     uint8     `db:"score_1"`
	Score2     uint8     `db:"score_2"`
	Rounds     uint8     `db:"rounds"`
	FinishedAt time.Time `db:"finished_at"`
}

type GameToCreate struct {
	ID     string
	Score1 uint8
	Score2 uint8
	Rounds uint8
}
