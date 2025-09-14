package session_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
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
						ID:          uuid.New(),
						TherapistID: uuid.New(),
						// SessionDate: time.Now(),
						// StartDateTime:   ptrString("10:00"),
						// EndDateTime:     ptrString("11:00"),
						Notes:     ptrString("Test session"),
						CreatedAt: ptrTime(time.Now()),
						UpdatedAt: ptrTime(time.Now()),
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
