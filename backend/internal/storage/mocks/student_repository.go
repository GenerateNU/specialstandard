package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockStudentRepository struct {
	mock.Mock
}

func (m *MockStudentRepository) GetStudents(ctx context.Context, pagination utils.Pagination) ([]models.Student, error) {
	args := m.Called(ctx, pagination)
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

func (m *MockStudentRepository) AddStudent(ctx context.Context, student models.Student) (models.Student, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(models.Student), args.Error(1)
}

func (m *MockStudentRepository) UpdateStudent(ctx context.Context, student models.Student) (models.Student, error) {
	args := m.Called(ctx, student)
	return args.Get(0).(models.Student), args.Error(1)
}

func (m *MockStudentRepository) GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination) ([]models.StudentSessionsOutput, error) {
	args := m.Called(ctx, studentID, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudentSessionsOutput), args.Error(1)
}

func (m *MockStudentRepository) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
