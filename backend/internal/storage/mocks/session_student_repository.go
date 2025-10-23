package mocks

import (
	"context"
	"specialstandard/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockSessionStudentRepository struct {
	mock.Mock
}

func (m *MockSessionStudentRepository) CreateSessionStudent(ctx context.Context, input *models.CreateSessionStudentInput) (*[]models.SessionStudent, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.SessionStudent), args.Error(1)
}

func (m *MockSessionStudentRepository) DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *MockSessionStudentRepository) PatchSessionStudent(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.SessionStudent), args.Error(1)
}
