package service_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/service/handler/session"
	"specialstandard/internal/utils"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/config"
	"specialstandard/internal/models"
	"specialstandard/internal/service"
	"specialstandard/internal/storage"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrString(s string) *string {
	return &s
}

// Health and Session Tests
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
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockSessionRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get sessions and default pagination",
			url:  "",
			mockSetup: func(m *mocks.MockSessionRepository) {
				sessions := []models.Session{
					{
						ID:            uuid.New(),
						TherapistID:   uuid.New(),
						StartDateTime: time.Now(),
						EndDateTime:   time.Now().Add(time.Hour),
						Notes:         ptrString("Test session"),
						CreatedAt:     ptrTime(time.Now()),
						UpdatedAt:     ptrTime(time.Now()),
					},
				}
				m.On("GetSessions", mock.Anything, utils.NewPagination()).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "/",
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("GetSessions", mock.Anything, utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name: "Default Pagination",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockSessionRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetSessions", mock.Anything, pagination).Return([]models.Session{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockSessionRepo)

			repo := &storage.Repository{
				Session: mockSessionRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("GET", "/api/v1/sessions"+tt.url, nil)
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}

// Student Integration Tests
func TestGetStudentsEndpoint(t *testing.T) {
	// Setup
	mockStudentRepo := new(mocks.MockStudentRepository)

	mockStudentRepo.On("GetStudents", mock.Anything, "", uuid.Nil, "", utils.NewPagination()).Return([]models.Student{
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

// Add these tests to your server_test.go file

func TestGetStudentsEndpoint_WithGradeFilter(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)

	expectedStudents := []models.Student{
		{
			ID:          uuid.New(),
			FirstName:   "John",
			LastName:    "Doe",
			DOB:         ptrTime(time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)),
			TherapistID: uuid.New(),
			Grade:       ptrString("5th"),
			IEP:         ptrString("Test IEP"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockStudentRepo.On("GetStudents", mock.Anything, "5th", uuid.Nil, "", mock.AnythingOfType("utils.Pagination")).Return(expectedStudents, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?grade=5th", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var students []models.Student
	json.NewDecoder(resp.Body).Decode(&students)
	assert.Len(t, students, 1)
	assert.Equal(t, "5th", *students[0].Grade)

	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentsEndpoint_WithTherapistFilter(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)
	therapistID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	expectedStudents := []models.Student{
		{
			ID:          uuid.New(),
			FirstName:   "Jane",
			LastName:    "Smith",
			DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
			TherapistID: therapistID,
			Grade:       ptrString("4th"),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockStudentRepo.On("GetStudents", mock.Anything, "", therapistID, "", mock.AnythingOfType("utils.Pagination")).Return(expectedStudents, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?therapist_id=123e4567-e89b-12d3-a456-426614174000", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var students []models.Student
	json.NewDecoder(resp.Body).Decode(&students)
	assert.Len(t, students, 1)
	assert.Equal(t, therapistID, students[0].TherapistID)

	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentsEndpoint_WithNameFilter(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)

	expectedStudents := []models.Student{
		{
			ID:        uuid.New(),
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			FirstName: "Johnny",
			LastName:  "Smith",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockStudentRepo.On("GetStudents", mock.Anything, "", uuid.Nil, "John", mock.AnythingOfType("utils.Pagination")).Return(expectedStudents, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?name=John", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var students []models.Student
	json.NewDecoder(resp.Body).Decode(&students)
	assert.Len(t, students, 2)

	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentsEndpoint_WithCombinedFilters(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)
	therapistID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")

	expectedStudents := []models.Student{
		{
			ID:          uuid.New(),
			FirstName:   "John",
			LastName:    "Doe",
			Grade:       ptrString("5th"),
			TherapistID: therapistID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockStudentRepo.On("GetStudents", mock.Anything, "5th", therapistID, "John", utils.Pagination{Page: 1, Limit: 5}).Return(expectedStudents, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?grade=5th&therapist_id=123e4567-e89b-12d3-a456-426614174000&name=John&page=1&limit=5", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var students []models.Student
	json.NewDecoder(resp.Body).Decode(&students)
	assert.Len(t, students, 1)
	assert.Equal(t, "5th", *students[0].Grade)
	assert.Equal(t, therapistID, students[0].TherapistID)

	mockStudentRepo.AssertExpectations(t)
}

func TestGetStudentsEndpoint_InvalidTherapistID(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?therapist_id=invalid-uuid", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Should not call repository if validation fails
	mockStudentRepo.AssertNotCalled(t, "GetStudents")
}

func TestGetStudentsEndpoint_EmptyFiltersIgnored(t *testing.T) {
	mockStudentRepo := new(mocks.MockStudentRepository)

	expectedStudents := []models.Student{
		{
			ID:        uuid.New(),
			FirstName: "Test",
			LastName:  "Student",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Empty string filters should be treated as no filter
	mockStudentRepo.On("GetStudents", mock.Anything, "", uuid.Nil, "", mock.AnythingOfType("utils.Pagination")).Return(expectedStudents, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
	}

	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/students?grade=&name=&therapist_id=", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	mockStudentRepo.AssertExpectations(t)
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
		ID:          uuid.New(),
		FirstName:   "John",
		LastName:    "Doe",
		Grade:       ptrString("5th"),
		TherapistID: uuid.New(),
		DOB:         ptrTime(time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)),
		IEP:         ptrString("Active IEP with speech therapy goals"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil)

	repo := &storage.Repository{
		Student: mockStudentRepo,
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
		ID:          studentID,
		FirstName:   "Emma",
		LastName:    "Johnson",
		Grade:       ptrString("5th"), // updated grade
		TherapistID: uuid.New(),
		DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
		IEP:         ptrString("Updated IEP with math accommodations"),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
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

func TestGetSessionByIDEndpoint(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		mockSetup      func(*mocks.MockSessionRepository)
		expectedStatus int
	}{
		{
			name:      "successful get session by valid UUID",
			sessionID: "123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func(m *mocks.MockSessionRepository) {
				session := models.Session{
					ID:          uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
					TherapistID: uuid.New(),
					Notes:       ptrString("Test session"),
				}
				m.On("GetSessionByID", mock.Anything, "123e4567-e89b-12d3-a456-426614174000").Return(&session, nil)
			},
			expectedStatus: 200,
		},
		{
			name:           "bad request for invalid UUID",
			sessionID:      "invalid-uuid",
			mockSetup:      func(m *mocks.MockSessionRepository) {}, // No mock calls expected
			expectedStatus: 400,
		},
		{
			name:      "session not found",
			sessionID: "123e4567-e89b-12d3-a456-426614174001", // Valid UUID but doesn't exist
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("GetSessionByID", mock.Anything, "123e4567-e89b-12d3-a456-426614174001").Return(nil, pgx.ErrNoRows)
			},
			expectedStatus: 404,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockSessionRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockSessionRepo)

			repo := &storage.Repository{
				Session: mockSessionRepo,
			}

			app := service.SetupApp(config.Config{}, repo)

			// Test
			req := httptest.NewRequest("GET", "/api/v1/sessions/"+tt.sessionID, nil)
			resp, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteSessionsEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		sessionID          uuid.UUID
		mockSetup          func(*mocks.MockSessionRepository, uuid.UUID)
		expectedStatusCode int
	}{
		{
			name:      "Success",
			sessionID: uuid.New(),
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSession", mock.Anything, id).Return("deleted", nil)
			},
			expectedStatusCode: 200,
		},
	}

	t.Run("Invalid UUID - Parsing Error", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionRepository)

		handler := session.NewHandler(mockRepo)
		app.Delete("/sessions/:id", handler.DeleteSessions)

		req := httptest.NewRequest("DELETE", "/sessions/0345", nil)
		res, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockSessionRepo, tt.sessionID)

			repo := &storage.Repository{
				Session: mockSessionRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("DELETE", "/api/v1/sessions/"+tt.sessionID.String(), nil)
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PostSessions(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionRepository)
		expectedStatusCode int
	}{
		{
			name: "Missing Items, Invalid JSON",
			payload: `{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z"
			}`,
			mockSetup:          func(m *mocks.MockSessionRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Empty Values that are actually Required",
			payload: `{
				"start_datetime": "",
				"end_datetime": "",
				"therapist_id": "00000000-0000-0000-0000-000000000000",
				"notes": ""
			}`,
			mockSetup:          func(m *mocks.MockSessionRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Foreign Key Violation: Therapist ID doesn't exist. DOCTOR WHO?!",
			payload: `{ 
				"start_datetime": "2025-09-14T10:00:00Z", 
				"end_datetime": "2025-09-14T11:00:00Z", 
				"therapist_id": "00000000-0000-0000-0000-000000000001", 
				"notes": "Test FK"
			}`,
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("PostSession", mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).Return(nil, errors.New("foreign key violation"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Start time and end time (check constraint violation)",
			payload: `{
				"start_datetime": "2025-09-14T11:00:00Z",
				"end_datetime": "2025-09-14T10:00:00Z",
				"therapist_id": "28eedfdc-81e1-44e5-a42c-022dc4c3b64d",
				"notes": "Check violation"
			}`,
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("PostSession", mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).Return(nil, errors.New("check constraint"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Success!",
			payload: `{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z",
				"therapist_id": "28eedfdc-81e1-44e5-a42c-022dc4c3b64d",
				"notes": "Test Session"
			}`,
			mockSetup: func(m *mocks.MockSessionRepository) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T10:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T11:00:00Z")
				sessionUUID := uuid.MustParse("28eedfdc-81e1-44e5-a42c-022dc4c3b64d")
				notes := ptrString("Test Session")
				now := time.Now()

				postSession := &models.PostSessionInput{
					StartTime:   startTime,
					EndTime:     endTime,
					TherapistID: sessionUUID,
					Notes:       notes,
				}

				session := &models.Session{
					ID:            uuid.New(),
					StartDateTime: startTime,
					EndDateTime:   endTime,
					TherapistID:   sessionUUID,
					Notes:         notes,
					CreatedAt:     &now,
					UpdatedAt:     &now,
				}
				m.On("PostSession", mock.Anything, postSession).Return(session, nil)
			},
			expectedStatusCode: fiber.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockSessionRepo)

			repo := &storage.Repository{
				Session: mockSessionRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("POST", "/api/v1/sessions", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockSessionRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchSessions(t *testing.T) {
	tests := []struct {
		id                 uuid.UUID
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionRepository, uuid.UUID)
		expectedStatusCode int
	}{
		{
			id:                 uuid.New(),
			name:               "Parsing PatchInputSession Error",
			payload:            `{"notes": "Missing quote}`,
			mockSetup:          func(m *mocks.MockSessionRepository, id uuid.UUID) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "Given ID not found",
			payload: `{"notes": "Trying to update non-existent"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				patch := &models.PatchSessionInput{
					Notes: ptrString("Trying to update non-existent"),
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, pgx.ErrNoRows)
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			id:      uuid.New(),
			name:    "foreign key violation - DOCTOR WHO? We don't know..",
			payload: `{"therapist_id": "00000000-0000-0000-0000-000000000999"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				therapistID := uuid.MustParse("00000000-0000-0000-0000-000000000999")
				patch := &models.PatchSessionInput{
					TherapistID: &therapistID,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, errors.New("foreign key"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "check constraint violation",
			payload: `{"start_datetime": "2025-09-14T14:00:00Z", "end_datetime": "2025-09-14T12:00:00Z"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T14:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				patch := &models.PatchSessionInput{
					StartTime: &startTime,
					EndTime:   &endTime,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, errors.New("check constraint"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "Successfully changed 1 field",
			payload: `{"notes": "The child seeks to be seen more than they wish to be heard"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				notes := ptrString("The child seeks to be seen more than they wish to be heard")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					Notes: notes,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: time.Now(),
					EndDateTime:   time.Now().Add(time.Hour),
					TherapistID:   uuid.New(),
					Notes:         notes,
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			id:      uuid.New(),
			name:    "Successfully changed multiple fields",
			payload: `{"start_datetime": "2025-09-14T12:00:00Z", "end_datetime": "2025-09-14T13:00:00Z"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T13:00:00Z")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					StartTime: &startTime,
					EndTime:   &endTime,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: startTime,
					EndDateTime:   endTime,
					TherapistID:   uuid.New(),
					Notes:         ptrString("Rescheduled for convenience"),
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			id:   uuid.New(),
			name: "Successfully changed all patchable fields",
			payload: `{
				"start_datetime": "2025-09-14T12:00:00Z", 
				"end_datetime": "2025-09-14T13:00:00Z", 
				"therapist_id": "28eedfdc-81e1-44e5-a42c-022dc4c3b64d", 
				"notes": "Starting Over"
			}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T13:00:00Z")
				therapistID := uuid.MustParse("28eedfdc-81e1-44e5-a42c-022dc4c3b64d")
				notes := ptrString("Starting Over")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					StartTime:   &startTime,
					EndTime:     &endTime,
					TherapistID: &therapistID,
					Notes:       notes,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: startTime,
					EndDateTime:   endTime,
					TherapistID:   therapistID,
					Notes:         notes,
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
	}

	t.Run("Bad UUID Request - Not a UUID", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionRepository)

		handler := session.NewHandler(mockRepo)
		app.Patch("/sessions/:id", handler.PatchSessions)

		req := httptest.NewRequest("PATCH", "/sessions/0345", nil)
		req.Header.Set("Content-Type", "application/json")
		res, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo, tt.id)

			handler := session.NewHandler(mockRepo)
			app.Patch("/sessions/:id", handler.PatchSessions)

			req := httptest.NewRequest("PATCH", "/sessions/"+tt.id.String(), strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetTherapistByIDEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(&models.Therapist{
		ID:        uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:     "matulakevin91@gmail.com",
		Active:    true,
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
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get therapists with default pagination",
			url:  "",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapists := []models.Therapist{
					{
						ID:        uuid.New(),
						FirstName: "Kevin",
						LastName:  "Matula",
						Email:     "matulakevin91@gmail.com",
						Active:    true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				m.On("GetTherapists", mock.Anything, utils.NewPagination()).Return(therapists, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("GetTherapists", mock.Anything, utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Pagination with parameters",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetTherapists", mock.Anything, pagination).Return([]models.Therapist{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTherapistRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockTherapistRepo)

			repo := &storage.Repository{
				Therapist: mockTherapistRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("GET", "/api/v1/therapists"+tt.url, nil)
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestCreateTherapistEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("CreateTherapist", mock.Anything).Return(&models.Therapist{
		ID:        uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:     "matulakevin91@gmail.com",
		Active:    true,
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

func TestDeleteTherapistEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("DeleteTherapist", mock.Anything, mock.AnythingOfType("string")).Return(nil)

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

func TestPatchTherapistEndpoint(t *testing.T) {
	// Setup
	mockTherapistRepo := new(mocks.MockTherapistRepository)

	mockTherapistRepo.On("PatchTherapist", mock.Anything, mock.AnythingOfType("string"), mock.Anything).Return(&models.Therapist{
		ID:        uuid.New(),
		FirstName: "Kevin",
		LastName:  "Matula",
		Email:     "matulakevin91@gmail.com",
		Active:    true,
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

func TestCreateResourceEndpoint(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	mockResourceRepo.On("CreateResource", mock.Anything, mock.Anything).Return(&models.Resource{
		ID:    uuid.New(),
		Title: ptrString("Resource1"),
		Type:  ptrString("doc"),
	}, nil)

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body := `{"title": "Resource1", "type": "doc"}`
	req := httptest.NewRequest("POST", "/api/v1/resources", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestGetResourcesEndpoint(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	mockResourceRepo.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, utils.NewPagination()).Return([]models.Resource{
		{
			ID:    uuid.New(),
			Title: ptrString("Resource1"),
			Type:  ptrString("doc"),
		},
	}, nil)

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/resources", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestGetResourceByIDEndpoint(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("GetResourceByID", mock.Anything, resourceID).Return(&models.Resource{
		ID:    resourceID,
		Title: ptrString("Resource1"),
		Type:  ptrString("doc"),
	}, nil)

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/resources/"+resourceID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestUpdateResourceEndpoint(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("UpdateResource", mock.Anything, resourceID, mock.Anything).Return(&models.Resource{
		ID:    resourceID,
		Title: ptrString("Updated Resource"),
		Type:  ptrString("doc"),
	}, nil)

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body := `{"title": "Updated Resource"}`
	req := httptest.NewRequest("PATCH", "/api/v1/resources/"+resourceID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestDeleteResourceEndpoint(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("DeleteResource", mock.Anything, resourceID).Return(nil)

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("DELETE", "/api/v1/resources/"+resourceID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestGetResourceByIDEndpoint_NotFound(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("GetResourceByID", mock.Anything, resourceID).Return((*models.Resource)(nil), errors.New("no rows in result set"))

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/resources/"+resourceID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestUpdateResourceEndpoint_NotFound(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("UpdateResource", mock.Anything, mock.Anything, mock.Anything).Return((*models.Resource)(nil), errors.New("no rows in result set"))

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body := `{"title": "Updated Resource"}`
	req := httptest.NewRequest("PATCH", "/api/v1/resources/"+resourceID.String(), strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestDeleteResourceEndpoint_NotFound(t *testing.T) {
	mockResourceRepo := new(mocks.MockResourceRepository)
	resourceID := uuid.New()
	mockResourceRepo.On("DeleteResource", mock.Anything, resourceID).Return(errors.New("resource not found"))

	repo := &storage.Repository{
		Resource: mockResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("DELETE", "/api/v1/resources/"+resourceID.String(), nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	mockResourceRepo.AssertExpectations(t)
}

func TestCreateSessionStudentEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionStudentRepository)
		expectedStatusCode int
	}{
		{
			name: "Successful creation",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true,
				"notes": "Student participated well in group activities"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				studentID := uuid.MustParse("987fcdeb-51a2-43d1-9c4f-123456789abc")

				m.On("CreateSessionStudent", mock.Anything, mock.AnythingOfType("*models.CreateSessionStudentInput")).Return(&models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   true,
					Notes:     ptrString("Student participated well in group activities"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatusCode: fiber.StatusCreated,
		},
		{
			name: "Missing session ID",
			payload: `{
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true
			}`,
			mockSetup:          func(m *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Invalid JSON format",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": /* missing comma */
			}`,
			mockSetup:          func(m *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Duplicate relationship",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true,
				"notes": "Duplicate entry"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("CreateSessionStudent", mock.Anything, mock.AnythingOfType("*models.CreateSessionStudentInput")).Return(nil, errors.New("duplicate key value violates unique constraint"))
			},
			expectedStatusCode: fiber.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionStudentRepo := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockSessionStudentRepo)

			repo := &storage.Repository{
				SessionStudent: mockSessionStudentRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("POST", "/api/v1/session_students", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockSessionStudentRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteSessionStudentEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionStudentRepository)
		expectedStatusCode int
	}{
		{
			name: "Successful deletion",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(nil)
			},
			expectedStatusCode: fiber.StatusNoContent,
		},
		{
			name: "Missing student ID",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000"
			}`,
			mockSetup:          func(m *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Invalid JSON format",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": /* missing comma */
			}`,
			mockSetup:          func(m *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Relationship not found",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(errors.New("no rows affected"))
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionStudentRepo := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockSessionStudentRepo)

			repo := &storage.Repository{
				SessionStudent: mockSessionStudentRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("DELETE", "/api/v1/session_students", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockSessionStudentRepo.AssertExpectations(t)
		})
	}
}

func TestPatchSessionStudentEndpoint(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionStudentRepository)
		expectedStatusCode int
	}{
		{
			name: "Successful patch - present only",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": false
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				studentID := uuid.MustParse("987fcdeb-51a2-43d1-9c4f-123456789abc")

				m.On("PatchSessionStudent", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(&models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   false,
					Notes:     ptrString("Original notes"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			name: "Successful patch - notes only",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"notes": "Updated progress notes"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				studentID := uuid.MustParse("987fcdeb-51a2-43d1-9c4f-123456789abc")

				m.On("PatchSessionStudent", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(&models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   true,
					Notes:     ptrString("Updated progress notes"),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			name: "Missing session ID",
			payload: `{
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true
			}`,
			mockSetup:          func(m *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Relationship not found",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("PatchSessionStudent", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(nil, errors.New("no rows affected"))
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			name: "Foreign key violation",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"student_id": "987fcdeb-51a2-43d1-9c4f-123456789abc",
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("PatchSessionStudent", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(nil, errors.New("foreign key violation"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSessionStudentRepo := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockSessionStudentRepo)

			repo := &storage.Repository{
				SessionStudent: mockSessionStudentRepo,
			}
			app := service.SetupApp(config.Config{}, repo)

			req := httptest.NewRequest("PATCH", "/api/v1/session_students", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, err := app.Test(req, -1)

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockSessionStudentRepo.AssertExpectations(t)
		})
	}
}

func TestGetResourcesBySessionIDEndpoint_Success(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)
	sessionID := uuid.New()

	expectedResources := []models.Resource{
		{
			ID:         uuid.New(),
			ThemeID:    uuid.New(),
			GradeLevel: ptrString("5th Grade"),
			Type:       ptrString("worksheet"),
			Title:      ptrString("Math Worksheet"),
			Category:   ptrString("math"),
			Content:    ptrString("Basic arithmetic"),
		},
		{
			ID:         uuid.New(),
			ThemeID:    uuid.New(),
			GradeLevel: ptrString("5th Grade"),
			Type:       ptrString("activity"),
			Title:      ptrString("Reading Activity"),
			Category:   ptrString("language"),
			Content:    ptrString("Comprehension exercise"),
		},
	}

	mockSessionResourceRepo.On("GetResourcesBySessionID", mock.Anything, sessionID, utils.NewPagination()).Return(expectedResources, nil)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/sessions/"+sessionID.String()+"/resources", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var resources []models.Resource
	err = json.NewDecoder(resp.Body).Decode(&resources)
	assert.NoError(t, err)
	assert.Len(t, resources, 2)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestGetResourcesBySessionIDEndpoint_EmptyArray(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)
	sessionID := uuid.New()

	mockSessionResourceRepo.On("GetResourcesBySessionID", mock.Anything, sessionID, utils.NewPagination()).Return([]models.Resource{}, nil)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/sessions/"+sessionID.String()+"/resources", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var resources []models.Resource
	err = json.NewDecoder(resp.Body).Decode(&resources)
	assert.NoError(t, err)
	assert.Empty(t, resources)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestGetResourcesBySessionIDEndpoint_InvalidUUID(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/sessions/invalid-uuid/resources", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockSessionResourceRepo.AssertNotCalled(t, "GetResourcesBySessionID")
}

func TestGetResourcesBySessionIDEndpoint_InternalError(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)
	sessionID := uuid.New()

	mockSessionResourceRepo.On("GetResourcesBySessionID", mock.Anything, sessionID, utils.NewPagination()).Return(nil, errors.New("database error"))

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("GET", "/api/v1/sessions/"+sessionID.String()+"/resources", nil)
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}

// POST /session-resource tests

func TestPostSessionResourceEndpoint_Success(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	createReq := models.CreateSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	expectedResponse := &models.SessionResource{
		SessionID:  createReq.SessionID,
		ResourceID: createReq.ResourceID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mockSessionResourceRepo.On("PostSessionResource", mock.Anything, createReq).Return(expectedResponse, nil)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 201, resp.StatusCode)

	var sessionResource models.SessionResource
	err = json.NewDecoder(resp.Body).Decode(&sessionResource)
	assert.NoError(t, err)
	assert.Equal(t, createReq.SessionID, sessionResource.SessionID)
	assert.Equal(t, createReq.ResourceID, sessionResource.ResourceID)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestPostSessionResourceEndpoint_SessionNotFound(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	createReq := models.CreateSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	pgErr := &pgconn.PgError{
		Code:   "23503",
		Detail: "Key (session_id)=(xxx) is not present in table",
	}

	mockSessionResourceRepo.On("PostSessionResource", mock.Anything, createReq).Return((*models.SessionResource)(nil), pgErr)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestPostSessionResourceEndpoint_ResourceNotFound(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	createReq := models.CreateSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	pgErr := &pgconn.PgError{
		Code:   "23503",
		Detail: "Key (resource_id)=(xxx) is not present in table",
	}

	mockSessionResourceRepo.On("PostSessionResource", mock.Anything, createReq).Return((*models.SessionResource)(nil), pgErr)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(createReq)
	req := httptest.NewRequest("POST", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestPostSessionResourceEndpoint_InvalidBody(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("POST", "/api/v1/session-resource", bytes.NewReader([]byte(`{"invalid": "json"`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockSessionResourceRepo.AssertNotCalled(t, "PostSessionResource")
}

func TestPostSessionResourceEndpoint_MissingFields(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body := `{"session_id": ""}`
	req := httptest.NewRequest("POST", "/api/v1/session-resource", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockSessionResourceRepo.AssertNotCalled(t, "PostSessionResource")
}

// DELETE /session-resource tests

func TestDeleteSessionResourceEndpoint_Success(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	deleteReq := models.DeleteSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	mockSessionResourceRepo.On("DeleteSessionResource", mock.Anything, deleteReq).Return(nil)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 204, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestDeleteSessionResourceEndpoint_NotFound(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	deleteReq := models.DeleteSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	mockSessionResourceRepo.On("DeleteSessionResource", mock.Anything, deleteReq).
		Return(fiber.NewError(fiber.StatusNotFound, "session resource relationship not found"))

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}

func TestDeleteSessionResourceEndpoint_InvalidBody(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	req := httptest.NewRequest("DELETE", "/api/v1/session-resource", bytes.NewReader([]byte(`{"invalid": "json"`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockSessionResourceRepo.AssertNotCalled(t, "DeleteSessionResource")
}

func TestDeleteSessionResourceEndpoint_MissingFields(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body := `{"session_id": ""}`
	req := httptest.NewRequest("DELETE", "/api/v1/session-resource", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	mockSessionResourceRepo.AssertNotCalled(t, "DeleteSessionResource")
}

func TestDeleteSessionResourceEndpoint_InternalError(t *testing.T) {
	mockSessionResourceRepo := new(mocks.MockSessionResourceRepository)

	deleteReq := models.DeleteSessionResource{
		SessionID:  uuid.New(),
		ResourceID: uuid.New(),
	}

	mockSessionResourceRepo.On("DeleteSessionResource", mock.Anything, deleteReq).Return(errors.New("database error"))

	repo := &storage.Repository{
		SessionResource: mockSessionResourceRepo,
	}
	app := service.SetupApp(config.Config{}, repo)

	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/api/v1/session-resource", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req, -1)

	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)

	mockSessionResourceRepo.AssertExpectations(t)
}
