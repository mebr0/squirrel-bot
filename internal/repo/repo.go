package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mebr0/squirrel-bot/internal/domain"
)

//go:generate mockgen -source=repo.go -destination=mocks/mock.go

type Players interface {
	// CreateOrUpdate write new user to database
	// If user already exists, update information with actual information
	CreateOrUpdate(ctx context.Context, player domain.Player) (domain.Player, error)
}

type Games interface {
	// List games of single player
	List(ctx context.Context, playerID int64) ([]domain.Game, error)
	// Create game results in database
	Create(ctx context.Context, toCreate domain.GameToCreate, players ...domain.PlayerTeam) (domain.Game, error)
}

type Repos struct {
	Players
	Games
}

func NewRepos(db *sqlx.DB) *Repos {
	return &Repos{
		Players: newPlayersRepo(db),
		Games:   newGamesRepo(db),
	}
}
