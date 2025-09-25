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

func (m *MockResourceRepository) GetResourceByID(ctx context.Context, id uuid.UUID) (*models.Resource, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Resource), args.Error(1)
}

func (m *MockResourceRepository) GetResources(ctx context.Context, theme_id uuid.UUID, gradeLevel, res_type, title, category, content string, date *time.Time, pagination utils.Pagination) ([]models.Resource, error) {
	args := m.Called(ctx, theme_id, gradeLevel, res_type, title, category, content, date, pagination)
	return args.Get(0).([]models.Resource), args.Error(1)
}

func (m *MockResourceRepository) UpdateResource(ctx context.Context, id uuid.UUID, resource models.UpdateResourceBody) (*models.Resource, error) {
	args := m.Called(ctx, id, resource)
	return args.Get(0).(*models.Resource), args.Error(1)
}
func (m *MockResourceRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
