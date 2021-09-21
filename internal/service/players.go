package service

import (
	"context"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/internal/repo"
)

type PlayersService struct {
	repo repo.Players
}

func newPlayersService(repo repo.Players) *PlayersService {
	return &PlayersService{
		repo: repo,
	}
}

func (s *PlayersService) Register(ctx context.Context, player domain.Player) (domain.Player, error) {
	return s.repo.CreateOrUpdate(ctx, player)
}
