package service_test

import (
	"net/http/httptest"
	"testing"
	"time"

	"specialstandard/internal/config"
	"specialstandard/internal/models"
	"specialstandard/internal/service"
	"specialstandard/internal/storage"
	"specialstandard/internal/storage/mocks"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string {
	return &s
}

func TestHealthEndpoint(t *testing.T) {
	// Setup
	mockSessionRepo := new(mocks.MockSessionRepository)
	repo := &storage.Repository{
		Session: mockSessionRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetSessionsEndpoint(t *testing.T) {
	// Setup
	mockSessionRepo := new(mocks.MockSessionRepository)

	mockSessionRepo.On("GetSessions", mock.Anything).Return([]models.Session{
		{
			ID:          uuid.New(),
			TherapistID: uuid.New(),
			Notes:       ptrString("Test session"),
		},
	}, nil)

	repo := &storage.Repository{
		Session: mockSessionRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/sessions", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetTherapistByIDEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(&models.Therapist{
		ID:         uuid.New(),
		First_name: "Kevin",
		Last_name:  "Matula",
		Email:      "matulakevin91@gmail.com",
		Active:     true,
		Created_at: time.Now(),
		Updated_at: time.Now(),
	}, nil)

	repo := &storage.Repository{
		Therapist: mockTherapistRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/therapists/9dad94d8-6534-4510-90d7-e4e97c175a65", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetTherapistsEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("GetTherapists", mock.Anything).Return([]models.Therapist{
		{
			ID:         uuid.New(),
			First_name: "Kevin",
			Last_name:  "Matula",
			Email:      "matulakevin91@gmail.com",
			Active:     true,
			Created_at: time.Now(),
			Updated_at: time.Now(),
		},
	}, nil)

	repo := &storage.Repository{
		Therapist: mockTherapistRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/therapists", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
