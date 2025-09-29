package mocks

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/stretchr/testify/mock"
)

type MockTherapistRepository struct {
	mock.Mock
}

func (m *MockTherapistRepository) GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error) {
	args := m.Called(ctx)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Therapist), args.Error(1)
}

func (m *MockTherapistRepository) GetTherapists(ctx context.Context, pagination utils.Pagination) ([]models.Therapist, error) {
	args := m.Called(ctx, pagination)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Therapist), args.Error(1)
}

func (m *MockTherapistRepository) CreateTherapist(ctx context.Context, therapist *models.CreateTherapistInput) (*models.Therapist, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Therapist), args.Error(1)
}

func (m *MockTherapistRepository) DeleteTherapist(ctx context.Context, therapistID string) error {
	args := m.Called(ctx, therapistID)
	return args.Error(0)
}

func (m *MockTherapistRepository) PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error) {
	args := m.Called(ctx, therapistID, updatedValue)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Therapist), args.Error(1)
}
