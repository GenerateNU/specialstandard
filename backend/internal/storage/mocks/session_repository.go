package mocks

import (
	"context"
	"specialstandard/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) GetSessions(ctx context.Context) ([]models.Session, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Session), args.Error(1)
}
