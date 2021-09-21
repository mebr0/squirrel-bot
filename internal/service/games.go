package service

import (
	"context"
	"github.com/mebr0/squirrel-bot/internal/domain"
	"github.com/mebr0/squirrel-bot/internal/repo"
)

type GamesService struct {
	repo repo.Games
}

func newGamesService(repo repo.Games) *GamesService {
	return &GamesService{
		repo: repo,
	}
}

func (s *GamesService) List(ctx context.Context, playerID int64) ([]domain.Game, error) {
	return s.repo.List(ctx, playerID)
}

func (s *GamesService) Save(ctx context.Context, toCreate domain.GameToCreate, players ...domain.PlayerTeam) (domain.Game, error) {
	return s.repo.Create(ctx, toCreate, players...)
}
