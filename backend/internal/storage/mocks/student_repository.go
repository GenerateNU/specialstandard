package mocks

import (
	"context"
	"specialstandard/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) GetStudents(ctx context.Context) ([]models.Student, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Student), args.Error(1)
}

func (m *MockStudentRepository) GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return models.Student{}, args.Error(1)
	}
	return args.Get(0).(models.Student), args.Error(1)
}

func (m *MockStudentRepository) AddStudent(ctx context.Context, student models.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentRepository) UpdateStudent(ctx context.Context, student models.Student) error {
	args := m.Called(ctx, student)
	return args.Error(0)
}

func (m *MockStudentRepository) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}