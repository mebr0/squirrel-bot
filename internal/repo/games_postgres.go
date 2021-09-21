package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mebr0/squirrel-bot/internal/domain"
)

type GamesRepo struct {
	db *sqlx.DB
}

func newGamesRepo(db *sqlx.DB) *GamesRepo {
	return &GamesRepo{
		db: db,
	}
}

func (r *GamesRepo) List(ctx context.Context, playerID int64) ([]domain.Game, error) {
	list := make([]domain.Game, 0)

	if err := r.db.SelectContext(ctx, &list, `
	SELECT g.id, g.score_1, g.score_2, g.rounds
	FROM games g 
    LEFT JOIN game_players gp on g.id = gp.game_id 
	LEFT JOIN players p on p.id = gp.player_id
	WHERE p.id = $1
	ORDER BY g.finished_at DESC`, playerID); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *GamesRepo) Create(ctx context.Context, toCreate domain.GameToCreate, players ...domain.PlayerTeam) (domain.Game, error) {
	tx, err := r.db.Begin()

	if err != nil {
		return domain.Game{}, err
	}

	row := tx.QueryRowContext(ctx,
		`INSERT INTO games(id, score_1, score_2, rounds) 
				VALUES ($1, $2, $3, $4) 
				RETURNING id, score_1, score_2, rounds`,
		toCreate.ID, toCreate.Score1, toCreate.Score2, toCreate.Rounds)

	var g domain.Game

	if err = row.Scan(&g.ID, &g.Score1, &g.Score2, &g.Rounds); err != nil {
		if err := tx.Rollback(); err != nil {
			return g, err
		}

		return g, err
	}

	for _, p := range players {
		if _, err = tx.ExecContext(ctx, "INSERT INTO game_players(game_id, player_id, team) VALUES ($1, $2, $3)",
			g.ID, p.ID, p.Team); err != nil {
			if err := tx.Rollback(); err != nil {
				return g, err
			}

			return g, err
		}
	}

	return g, tx.Commit()
}
