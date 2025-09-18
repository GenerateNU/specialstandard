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
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)


func ptrTime(t time.Time) *time.Time {
    return &t
}

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
			DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
			TherapistID: uuid.New(),
			Grade:       ptrString("4th"),
			IEP:         ptrString("Reading comprehension support"),
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
		DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
		TherapistID: uuid.New(),
		Grade:       ptrString("4th"),
		IEP:         ptrString("Reading comprehension support"),
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
	mockStudentRepo.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
    ID: uuid.New(),
    FirstName: "John",
    LastName: "Doe", 
    Grade:       ptrString("5th"),
    TherapistID: uuid.New(),
    DOB:         ptrTime(time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)),
    IEP: ptrString("Active IEP with speech therapy goals"),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
		// TODO: Add Therapist mock when Kevin's therapist repository is merged
		// Currently therapist validation is commented out in AddStudent handler
	}

	app := service.SetupApp(config.Config{}, repo)

	testTherapistID := uuid.New()
	
	body := fmt.Sprintf(`{
		"first_name": "John",
		"last_name": "Doe",
		"dob": "2010-05-15",
		"therapist_id": "%s",
		"grade": "5th",
		"iep": "Active IEP with speech therapy goals"
	}`, testTherapistID.String())

	req := httptest.NewRequest("POST", "/api/v1/students", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	mockStudentRepo.AssertExpectations(t)
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
		DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
		TherapistID: uuid.New(),
		Grade:       ptrString("4th"),
		IEP:         ptrString("Original IEP"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	// Mock UpdateStudent call
	mockStudentRepo.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
    ID: studentID,
    FirstName: "Emma", 
    LastName: "Johnson",
    Grade:       ptrString("5th"), // updated grade
    TherapistID: uuid.New(),
    DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
    IEP: ptrString("Updated IEP with math accommodations"),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}, nil)

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