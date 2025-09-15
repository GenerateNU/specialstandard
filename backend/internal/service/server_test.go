package service_test

import (
	"net/http/httptest"
	"testing"

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
			name:      "Invalid ID - Not Found / Doesn't Exist",
			sessionID: uuid.New(),
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSessions", mock.Anything, id).Return("", nil)
			},
			expectedStatusCode: 404,
		},
		{
			name:      "Success",
			sessionID: uuid.New(),
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSessions", mock.Anything, id).Return("deleted", nil)
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
