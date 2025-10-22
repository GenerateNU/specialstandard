package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockSessionRepository struct {
	mock.Mock
}

func (m *MockSessionRepository) GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest) ([]models.Session, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Session), args.Error(1)
}

func (m *MockSessionRepository) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockSessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *MockSessionRepository) PostSession(ctx context.Context, session *models.PostSessionInput) (*[]models.Session, error) {
	args := m.Called(ctx, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.Session), args.Error(1)
}

func (m *MockSessionRepository) PatchSession(ctx context.Context, id uuid.UUID, session *models.PatchSessionInput) (*models.Session, error) {
	args := m.Called(ctx, id, session)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Session), args.Error(1)
}

func (m *MockSessionRepository) GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.SessionStudentsOutput, error) {
	args := m.Called(ctx, sessionID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SessionStudentsOutput), args.Error(1)
}
