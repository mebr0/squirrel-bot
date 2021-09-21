package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/mebr0/squirrel-bot/internal/domain"
	mockRepo "github.com/mebr0/squirrel-bot/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	playerID int64 = 1
)

func mockGamesService(t *testing.T) (*GamesService, *mockRepo.MockGames) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	gRepo := mockRepo.NewMockGames(mockCtl)

	s := newGamesService(gRepo)

	return s, gRepo
}

func TestGamesService_List(t *testing.T) {
	s, gRepo := mockGamesService(t)

	ctx := context.Background()

	gRepo.EXPECT().List(ctx, playerID).Return([]domain.Game{}, nil)

	g, err := s.List(ctx, playerID)

	require.NoError(t, err)
	require.IsType(t, []domain.Game{}, g)
}

func TestGamesService_Save(t *testing.T) {
	s, gRepo := mockGamesService(t)

	ctx := context.Background()

	gRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.Game{}, nil)

	g, err := s.Save(ctx, domain.GameToCreate{})

	require.NoError(t, err)
	require.IsType(t, domain.Game{}, g)
}
