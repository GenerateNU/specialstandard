package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockResourceRepository struct {
	mock.Mock
}

func (m *MockResourceRepository) CreateResource(ctx context.Context, body models.ResourceBody) (*models.Resource, error) {
	args := m.Called(ctx, body)
	return args.Get(0).(*models.Resource), args.Error(1)
}

func (m *MockResourceRepository) GetResourceByID(ctx context.Context, id uuid.UUID) (*models.ResourceWithTheme, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.ResourceWithTheme), args.Error(1)
}

func (m *MockResourceRepository) GetResources(ctx context.Context, themeID uuid.UUID, gradeLevel, resType, title, category, content, themeName string, date *time.Time, themeMonth, themeYear int, pagination utils.Pagination) ([]models.ResourceWithTheme, error) {
	args := m.Called(ctx, themeID, gradeLevel, resType, title, category, content, date, themeName, themeMonth, themeYear, pagination)
	return args.Get(0).([]models.ResourceWithTheme), args.Error(1)
}

func (m *MockResourceRepository) UpdateResource(ctx context.Context, id uuid.UUID, resource models.UpdateResourceBody) (*models.Resource, error) {
	args := m.Called(ctx, id, resource)
	return args.Get(0).(*models.Resource), args.Error(1)
}
func (m *MockResourceRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
