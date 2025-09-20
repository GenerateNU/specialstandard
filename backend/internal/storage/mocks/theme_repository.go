package mocks

import (
	"context"
	"specialstandard/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockThemeRepository struct {
	mock.Mock
}

func (m *MockThemeRepository) CreateTheme(ctx context.Context, theme *models.CreateThemeInput) (*models.Theme, error) {
	args := m.Called(ctx, theme)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Theme), args.Error(1)
}

func (m *MockThemeRepository) GetThemes(ctx context.Context) ([]models.Theme, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Theme), args.Error(1)
}

func (m *MockThemeRepository) GetThemeByID(ctx context.Context, id uuid.UUID) (*models.Theme, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Theme), args.Error(1)
}

func (m *MockThemeRepository) UpdateTheme(ctx context.Context, id uuid.UUID, theme *models.UpdateThemeInput) (*models.Theme, error) {
	args := m.Called(ctx, id, theme)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Theme), args.Error(1)
}

func (m *MockThemeRepository) DeleteTheme(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
