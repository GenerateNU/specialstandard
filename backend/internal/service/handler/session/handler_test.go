package session_test

import (
	"errors"
	"net/http/httptest"
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
						SessionDate: time.Now(),
						StartTime:   ptrString("10:00"),
						EndTime:     ptrString("11:00"),
						Notes:       ptrString("Test session"),
						CreatedAt:   ptrTime(time.Now()),
						UpdatedAt:   ptrTime(time.Now()),
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
