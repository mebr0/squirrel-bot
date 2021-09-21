package service

import (
	"context"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/internal/repo"
)

type Players interface {
	// Register create new player in system or update information about him
	// Because users can change their user and full names
	Register(ctx context.Context, player domain.Player) (domain.Player, error)
}

type Services struct {
	Players
}

func NewServices(repos *repo.Repos) *Services {
	return &Services{
		Players: newPlayersService(repos.Players),
	}
}
