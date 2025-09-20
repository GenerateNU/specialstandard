package theme_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/theme"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestHandler_CreateTheme(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful create theme",
			mockSetup: func(m *mocks.MockThemeRepository) {
				theme := &models.Theme{
					ID:        uuid.New(),
					Name:      "Spring",
					Month:     3,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(theme, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Post("/themes", handler.CreateTheme)

			body := `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			}`

			req := httptest.NewRequest("POST", "/themes", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetThemes(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get themes",
			mockSetup: func(m *mocks.MockThemeRepository) {
				themes := []models.Theme{
					{
						ID:        uuid.New(),
						Name:      "Spring",
						Month:     3,
						Year:      2024,
						CreatedAt: ptrTime(time.Now()),
						UpdatedAt: ptrTime(time.Now()),
					},
				}
				m.On("GetThemes", mock.Anything).Return(themes, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemes", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Get("/themes", handler.GetThemes)

			req := httptest.NewRequest("GET", "/themes", nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetThemeByID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get theme by id",
			mockSetup: func(m *mocks.MockThemeRepository) {
				theme := &models.Theme{
					ID:        uuid.New(),
					Name:      "Spring",
					Month:     3,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}
				m.On("GetThemeByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(theme, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "theme not found",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemeByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, pgx.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemeByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Get("/themes/:id", handler.GetThemeByID)

			req := httptest.NewRequest("GET", "/themes/"+uuid.New().String(), nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchTheme(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful patch theme",
			mockSetup: func(m *mocks.MockThemeRepository) {
				theme := &models.Theme{
					ID:        uuid.New(),
					Name:      "Updated Spring",
					Month:     4,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}
				m.On("UpdateTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(theme, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "theme not found",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, pgx.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Patch("/themes/:id", handler.PatchTheme)

			body := `{
				"name": "Updated Spring",
				"month": 4
			}`

			req := httptest.NewRequest("PATCH", "/themes/"+uuid.New().String(), strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteTheme(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful delete theme",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "theme not found",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(errors.New("theme not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Delete("/themes/:id", handler.DeleteTheme)

			req := httptest.NewRequest("DELETE", "/themes/"+uuid.New().String(), nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}