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

func (m *MockStudentRepository) GetStudents(ctx context.Context, grade *int, therapistID uuid.UUID, name string, pagination utils.Pagination) ([]models.Student, error) {
	args := m.Called(ctx, grade, therapistID, name, pagination)
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

func (m *MockStudentRepository) GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRepositoryRequest) ([]models.StudentSessionsOutput, error) {
	args := m.Called(ctx, studentID, pagination, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudentSessionsOutput), args.Error(1)
}

func (m *MockStudentRepository) GetStudentRatings(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRatingsRequest) ([]models.StudentSessionsWithRatingsOutput, error) {
	args := m.Called(ctx, studentID, pagination, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.StudentSessionsWithRatingsOutput), args.Error(1)
}

func (m *MockStudentRepository) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStudentRepository) PromoteStudents(ctx context.Context, input models.PromoteStudentsInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}
