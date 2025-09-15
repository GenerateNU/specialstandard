package service_test

import (
	"net/http/httptest"
	"testing"

	"specialstandard/internal/config"
	"specialstandard/internal/models"
	"specialstandard/internal/service"
	"specialstandard/internal/storage"
	"specialstandard/internal/storage/mocks"
	"strings"
	"time"

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


// Student Integration Tests

func TestGetStudentsEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)

	mockStudentRepo.On("GetStudents", mock.Anything).Return([]models.Student{
		{
			ID:          uuid.New(),
			FirstName:   "Emma",
			LastName:    "Johnson",
			DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
			TherapistID: uuid.New(),
			Grade:       "4th",
			IEP:         "Reading comprehension support",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/students", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetStudentByIDEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)
	studentID := uuid.New()

	mockStudentRepo.On("GetStudent", mock.Anything, studentID).Return(models.Student{
		ID:          studentID,
		FirstName:   "Emma",
		LastName:    "Johnson",
		DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
		TherapistID: uuid.New(),
		Grade:       "4th",
		IEP:         "Reading comprehension support",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("GET", "/api/v1/students/"+studentID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateStudentEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)

	mockStudentRepo.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	body := `{
		"first_name": "John",
		"last_name": "Doe",
		"dob": "2010-05-15",
		"therapist_id": "9dad94d8-6534-4510-90d7-e4e97c175a65",
		"grade": "5th",
		"iep": "Active IEP with speech therapy goals"
	}`

	req := httptest.NewRequest("POST", "/api/v1/students", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
}

func TestUpdateStudentEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)
	studentID := uuid.New()

	// Mock GetStudent call (UpdateStudent handler calls this first)
	mockStudentRepo.On("GetStudent", mock.Anything, studentID).Return(models.Student{
		ID:          studentID,
		FirstName:   "Emma",
		LastName:    "Johnson",
		DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
		TherapistID: uuid.New(),
		Grade:       "4th",
		IEP:         "Original IEP",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Mock UpdateStudent call
	mockStudentRepo.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	body := `{
		"grade": "5th",
		"iep": "Updated IEP with math accommodations"
	}`

	req := httptest.NewRequest("PATCH", "/api/v1/students/"+studentID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestDeleteStudentEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)
	studentID := uuid.New()

	mockStudentRepo.On("DeleteStudent", mock.Anything, studentID).Return(nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	// Test
	req := httptest.NewRequest("DELETE", "/api/v1/students/"+studentID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
}