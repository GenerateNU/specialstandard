package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/dbinterface"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/mock"
)

type MockSessionStudentRepository struct {
	mock.Mock
}

// CreateSessionStudent - Returns *[]models.SessionStudent to match the interface
func (m *MockSessionStudentRepository) CreateSessionStudent(ctx context.Context, q dbinterface.Queryable, input *models.CreateSessionStudentInput) (*[]models.SessionStudent, error) {
	args := m.Called(ctx, q, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.SessionStudent), args.Error(1)
}

// DeleteSessionStudent - Updated to accept queryable as second parameter for transaction support
func (m *MockSessionStudentRepository) DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// PatchSessionStudent - This method seems to be missing, adding it based on common patterns
func (m *MockSessionStudentRepository) PatchSessionStudent(ctx context.Context, input *models.PatchSessionStudentInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// RateStudentSession - Updated to include queryable parameter
func (m *MockSessionStudentRepository) RateStudentSession(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, []models.SessionRating, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, nil, args.Error(2)
	}
	return args.Get(0).(*models.SessionStudent), args.Get(1).([]models.SessionRating), args.Error(2)
}

// GetDB returns the database pool
func (m *MockSessionStudentRepository) GetDB() *pgxpool.Pool {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*pgxpool.Pool)
}
