package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/dbinterface"
	"specialstandard/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockSessionRepository struct {
	mock.Mock
}

// Update from 3 parameters to 4 parameters
func (m *MockSessionRepository) GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest, therapistID uuid.UUID) ([]models.Session, error) {
	args := m.Called(ctx, pagination, filter, therapistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Session), args.Error(1)
}

// Also update GetSessionStudents to include therapistID
func (m *MockSessionRepository) GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination, therapistID uuid.UUID) ([]models.SessionStudentsOutput, error) {
	args := m.Called(ctx, sessionID, pagination, therapistID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SessionStudentsOutput), args.Error(1)
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
	return args.Error(0)
}

func (m *MockSessionRepository) PostSession(ctx context.Context, q dbinterface.Queryable, session *models.PostSessionInput) (*[]models.Session, error) {
	args := m.Called(ctx, q, session)
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

func (m *MockSessionRepository) GetDB() *pgxpool.Pool {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgxpool.Pool)
}
