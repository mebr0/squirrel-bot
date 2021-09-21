package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mebr0/squirrel-bot/internal/domain"
)

type PlayersRepo struct {
	db *sqlx.DB
}

func newPlayersRepo(db *sqlx.DB) *PlayersRepo {
	return &PlayersRepo{
		db: db,
	}
}

func (r *PlayersRepo) CreateOrUpdate(ctx context.Context, player domain.Player) (domain.Player, error) {
	row := r.db.QueryRowContext(ctx,
		`INSERT INTO players(id, username, first_name, last_name) 
				VALUES ($1, $2, $3, $4) 
				ON CONFLICT (id) 
				DO UPDATE SET username=$2, first_name=$3, last_name=$4
				RETURNING id, username, first_name, last_name`,
		player.ID, player.Username, player.FirstName, player.LastName)

	var p domain.Player

	if err := row.Scan(&p.ID, &p.Username, &p.FirstName, &p.LastName); err != nil {
		return p, err
	}

	return p, nil
}
