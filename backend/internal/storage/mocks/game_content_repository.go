package mocks

import (
	"context"
	"specialstandard/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockGameContentRepository struct {
	mock.Mock
}

func (m *MockGameContentRepository) GetGameContents(ctx context.Context, req models.GetGameContentRequest) ([]models.GameContent, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.GameContent), args.Error(1)
}
