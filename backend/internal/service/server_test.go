package service_test

import (
	"errors"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/service/handler/session"
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
