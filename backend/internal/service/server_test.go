package service_test

import (
	"net/http/httptest"
	"strings"
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
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:      "matulakevin91@gmail.com",
		Active:     true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
			FirstName: "Kevin",
			LastName:  "Matula",
			Email:      "matulakevin91@gmail.com",
			Active:     true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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

func TestCreateTherapistEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("CreateTherapist", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(&models.Therapist{
		ID:         uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:      "matulakevin91@gmail.com",
		Active:     true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	repo := &storage.Repository{
		Therapist: mockTherapistRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	body := `{
    "first_name": "Kevin",
    "last_name": "Matula",
    "email": "matulakevin91@gmail.com"
	}`

	req := httptest.NewRequest("POST", "/api/v1/therapists", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestDeleteTherapist(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("DeleteTherapist", mock.Anything, mock.AnythingOfType("string")).Return(&models.Therapist{
		ID:         uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:      "matulakevin91@gmail.com",
		Active:     true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	repo := &storage.Repository{
		Therapist: mockTherapistRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("DELETE", "/api/v1/therapists/4a9a4e58-ea6c-496a-915f-3e8214e77112", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestPatchTherapist(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("PatchTherapist", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(&models.Therapist{
		ID:         uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:      "matulakevin91@gmail.com",
		Active:     true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil)

	repo := &storage.Repository{
		Therapist: mockTherapistRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	body := `{
    "first_name": "Kevin",
    "last_name": "Matula",
    "email": "matulakevin91@gmail.com"
	}`

	// Test
	req := httptest.NewRequest("PATCH", "/api/v1/therapists/4a9a4e58-ea6c-496a-915f-3e8214e77112", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
