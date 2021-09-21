package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/mebr0/squirrel-bot/internal/domain"
	mockRepo "github.com/mebr0/squirrel-bot/internal/repo/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func mockPlayersService(t *testing.T) (*PlayersService, *mockRepo.MockPlayers) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	pRepo := mockRepo.NewMockPlayers(mockCtl)

	s := newPlayersService(pRepo)

	return s, pRepo
}

func TestPlayersService_Register(t *testing.T) {
	s, pRepo := mockPlayersService(t)

	ctx := context.Background()

	pRepo.EXPECT().CreateOrUpdate(ctx, gomock.Any()).Return(domain.Player{}, nil)

	p, err := s.Register(ctx, domain.Player{})

	require.NoError(t, err)
	require.IsType(t, domain.Player{}, p)
}
