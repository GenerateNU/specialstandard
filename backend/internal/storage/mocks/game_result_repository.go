package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/stretchr/testify/mock"
)

type MockGameResultRepository struct {
	mock.Mock
}

func (m *MockGameResultRepository) GetGameResults(ctx context.Context, inputQuery *models.GetGameResultQuery, pagination utils.Pagination) ([]models.GameResult, error) {
	args := m.Called(ctx, query, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.GameResult), args.Error(1)
}

func (m *MockGameResultRepository) PostGameResult(ctx context.Context, input models.PostGameResult) (*models.GameResult, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.GameResult), args.Error(1)
}
