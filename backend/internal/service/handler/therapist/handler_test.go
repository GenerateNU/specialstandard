package therapist_test

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/therapist"
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

func TestHandler_GetTherapistByID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get therapist by id",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapist := &models.Therapist{
						ID:          uuid.New(),
						First_name:  "Kevin",
						Last_name:   "Matula",
						Email:       "matulakevin91@gmail.com",
						Active:      true,
						Created_at:   time.Now(),
						Updated_at:   time.Now(),
					}
				m.On("GetTherapistByID", mock.Anything).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("GetTherapistByID", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Get("/therapists/:id", handler.GetTherapistByID)

			// Make request
			req := httptest.NewRequest("GET", "/therapists/9dad94d8-6534-4510-90d7-e4e97c175a65", nil)
			resp, _ := app.Test(req, -1)

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
