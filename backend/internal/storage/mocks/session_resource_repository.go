package mocks

import (
	"context"
	"specialstandard/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSessionResourceRepository struct {
	mock.Mock
}

func (m *MockSessionResourceRepository) PostSessionResource(ctx context.Context, sessionResource models.CreateSessionResource) (*models.SessionResource, error) {
	args := m.Called(mock.Anything, sessionResource)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SessionResource), args.Error(1)
}

func (m *MockSessionResourceRepository) DeleteSessionResource(ctx context.Context, sessionResource models.DeleteSessionResource) error {
	args := m.Called(mock.Anything, sessionResource)
	return args.Error(0)
}

func (m *MockSessionResourceRepository) GetResourcesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]models.Resource, error) {
	args := m.Called(mock.Anything, sessionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Resource), args.Error(1)
}
