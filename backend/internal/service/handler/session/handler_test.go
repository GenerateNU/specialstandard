package session_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/session"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestHandler_GetSessions(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockSessionRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get sessions",
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
				m.On("GetSessions", mock.Anything).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("GetSessions", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo)

			handler := session.NewHandler(mockRepo)
			app.Get("/sessions", handler.GetSessions)

			// Make request
			req := httptest.NewRequest("GET", "/sessions", nil)
			resp, _ := app.Test(req, -1)

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteSessions(t *testing.T) {
	tests := []struct {
		id             uuid.UUID
		name           string
		mockSetup      func(*mocks.MockSessionRepository, uuid.UUID)
		expectedStatus int
		wantErr        bool
	}{
		{
			id:   uuid.New(),
			name: "Successful Delete Session",
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSessions", mock.Anything, id).Return("deleted", nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			id:   uuid.New(),
			name: "internal server error",
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSessions", mock.Anything, id).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	t.Run("Bad UUID Request", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionRepository)

		handler := session.NewHandler(mockRepo)
		app.Delete("/sessions/:id", handler.DeleteSessions)

		req := httptest.NewRequest("DELETE", "/sessions/1234", nil)
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
			app.Delete("/sessions/:id", handler.DeleteSessions)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/sessions/%s", tt.id.String()), nil)
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockRepo.AssertExpectations(t)
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
			name: "Empty Values that are Required",
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
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T10:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T11:00:00Z")

				session := &models.PostSessionInput{
					StartTime:   startTime,
					EndTime:     endTime,
					TherapistID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
					Notes:       ptrString("Test FK"),
				}
				m.On("PostSessions", mock.Anything, session).Return(nil, errors.New("foreign key violation"))
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
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T11:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T10:00:00Z")

				session := &models.PostSessionInput{
					StartTime:   startTime,
					EndTime:     endTime,
					TherapistID: uuid.MustParse("28eedfdc-81e1-44e5-a42c-022dc4c3b64d"),
					Notes:       ptrString("Check violation"),
				}
				m.On("PostSessions", mock.Anything, session).Return(nil, errors.New("check constraint"))
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
				m.On("PostSessions", mock.Anything, postSession).Return(session, nil)
			},
			expectedStatusCode: fiber.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo)

			handler := session.NewHandler(mockRepo)
			app.Post("/sessions", handler.PostSessions)

			req := httptest.NewRequest("POST", "/sessions", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
