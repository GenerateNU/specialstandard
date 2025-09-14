package mocks

import (
	"context"
	"specialstandard/internal/models"

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

func (m *MockTherapistRepository) GetTherapists(ctx context.Context) ([]models.Therapist, error) {
	args := m.Called(ctx)
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

func (m *MockTherapistRepository) DeleteTherapist(ctx context.Context, therapistID string) (string, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return "User Deleted Successfully", args.Error(1)
}
