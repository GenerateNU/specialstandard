package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

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

func (m *MockSessionResourceRepository) GetResourcesBySessionID(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.Resource, error) {
	args := m.Called(ctx, sessionID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Resource), args.Error(1)
}
