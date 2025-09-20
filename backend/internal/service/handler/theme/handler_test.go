package theme_test

import (
	"encoding/json"
	"errors"
	"io"
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestHandler_CreateTheme(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful create theme",
			requestBody: `{
				"name": "Winter Wonderland",
				"month": 12,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(&models.Theme{
					ID:        uuid.New(),
					Name:      "Winter Wonderland",
					Month:     12,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "successful create theme spring",
			requestBody: `{
				"name": "Spring Blossoms",
				"month": 3,
				"year": 2025
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(&models.Theme{
					ID:        uuid.New(),
					Name:      "Spring Blossoms",
					Month:     3,
					Year:      2025,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name:        "invalid JSON body",
			requestBody: `{"name": "Winter" /* missing comma */}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing required fields",
			requestBody: `{"name": "Winter"}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid month too low",
			requestBody: `{
				"name": "Invalid Month",
				"month": 0,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid month too high",
			requestBody: `{
				"name": "Invalid Month",
				"month": 13,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid year too low",
			requestBody: `{
				"name": "Invalid Year",
				"month": 6,
				"year": 1899
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid year too high",
			requestBody: `{
				"name": "Invalid Year",
				"month": 6,
				"year": 2101
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "repository error",
			requestBody: `{
				"name": "Test Theme",
				"month": 6,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Post("/themes", handler.CreateTheme)

			// Make request
			req := httptest.NewRequest("POST", "/themes", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusCreated {
				// Success case - validate created theme data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var theme models.Theme
				err = json.Unmarshal(body, &theme)
				assert.NoError(t, err)

				// Validate response data
				switch tt.name {
				case "successful create theme":
					assert.Equal(t, "Winter Wonderland", theme.Name)
					assert.Equal(t, 12, theme.Month)
					assert.Equal(t, 2024, theme.Year)
				case "successful create theme spring":
					assert.Equal(t, "Spring Blossoms", theme.Name)
					assert.Equal(t, 3, theme.Month)
					assert.Equal(t, 2025, theme.Year)
				}

				// Validate that UUID was generated
				assert.NotEqual(t, uuid.Nil, theme.ID)
				assert.NotNil(t, theme.CreatedAt)
				assert.NotNil(t, theme.UpdatedAt)
			}
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
						Name:      "Winter Wonderland",
						Month:     12,
						Year:      2024,
						CreatedAt: ptrTime(time.Now()),
						UpdatedAt: ptrTime(time.Now()),
					},
					{
						ID:        uuid.New(),
						Name:      "Spring Blossoms",
						Month:     3,
						Year:      2025,
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
			name: "empty themes list",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemes", mock.Anything).Return([]models.Theme{}, nil)
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

			// Response body validation
			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var themes []models.Theme
				err = json.Unmarshal(body, &themes)
				assert.NoError(t, err)

				switch tt.name {
				case "empty themes list":
					assert.Len(t, themes, 0)
				case "successful get themes":
					assert.Len(t, themes, 2)
					assert.Equal(t, "Winter Wonderland", themes[0].Name)
					assert.Equal(t, "Spring Blossoms", themes[1].Name)
				}
			}
		})
	}
}

func TestHandler_GetThemeByID(t *testing.T) {
	themeID := uuid.New()

	tests := []struct {
		name           string
		themeID        string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:    "successful get theme",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				theme := &models.Theme{
					ID:        themeID,
					Name:      "Winter Wonderland",
					Month:     12,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}
				m.On("GetThemeByID", mock.Anything, themeID).Return(theme, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:    "theme not found",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemeByID", mock.Anything, themeID).Return(nil, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:    "invalid UUID format",
			themeID: "invalid-uuid",
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:    "repository error",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemeByID", mock.Anything, themeID).Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Get("/themes/:id", handler.GetThemeByID)

			// Make request
			req := httptest.NewRequest("GET", "/themes/"+tt.themeID, nil)
			resp, _ := app.Test(req, -1)

			// Basic assertions
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				// Success case - validate theme data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var theme models.Theme
				err = json.Unmarshal(body, &theme)
				assert.NoError(t, err)

				// Validate the theme data
				assert.Equal(t, "Winter Wonderland", theme.Name)
				assert.Equal(t, 12, theme.Month)
				assert.Equal(t, 2024, theme.Year)
				assert.Equal(t, themeID, theme.ID)
			}
		})
	}
}

func TestHandler_PatchTheme(t *testing.T) {
	themeID := uuid.New()

	tests := []struct {
		name           string
		themeID        string
		requestBody    string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:        "update name only",
			themeID:     themeID.String(),
			requestBody: `{"name": "Updated Winter Theme"}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, themeID, mock.AnythingOfType("*models.UpdateThemeInput")).Return(&models.Theme{
					ID:        themeID,
					Name:      "Updated Winter Theme",
					Month:     12,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "update month and year",
			themeID:     themeID.String(),
			requestBody: `{"month": 6, "year": 2025}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, themeID, mock.AnythingOfType("*models.UpdateThemeInput")).Return(&models.Theme{
					ID:        themeID,
					Name:      "Winter Wonderland",
					Month:     6,
					Year:      2025,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "update all fields",
			themeID:     themeID.String(),
			requestBody: `{"name": "Summer Fun", "month": 7, "year": 2025}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, themeID, mock.AnythingOfType("*models.UpdateThemeInput")).Return(&models.Theme{
					ID:        themeID,
					Name:      "Summer Fun",
					Month:     7,
					Year:      2025,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "invalid UUID format",
			themeID:     "invalid-uuid",
			requestBody: `{"name": "Test"}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "theme not found",
			themeID:     themeID.String(),
			requestBody: `{"name": "Test"}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, themeID, mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:        "invalid JSON body",
			themeID:     themeID.String(),
			requestBody: `{"name": "Test" /* missing comma */}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid month in update",
			themeID:     themeID.String(),
			requestBody: `{"month": 13}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "repository error",
			themeID:     themeID.String(),
			requestBody: `{"name": "Test"}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("UpdateTheme", mock.Anything, themeID, mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockThemeRepository)
			tt.mockSetup(mockRepo)

			handler := theme.NewHandler(mockRepo)
			app.Patch("/themes/:id", handler.PatchTheme)

			// Make request
			url := "/themes/" + tt.themeID
			req := httptest.NewRequest("PATCH", url, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				// Success case - validate updated theme data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var theme models.Theme
				err = json.Unmarshal(body, &theme)
				assert.NoError(t, err)

				assert.Equal(t, themeID, theme.ID)
				assert.NotNil(t, theme.CreatedAt)
				assert.NotNil(t, theme.UpdatedAt)

				// Validate specific updates
				switch tt.name {
				case "update name only":
					assert.Equal(t, "Updated Winter Theme", theme.Name)
				case "update month and year":
					assert.Equal(t, 6, theme.Month)
					assert.Equal(t, 2025, theme.Year)
				case "update all fields":
					assert.Equal(t, "Summer Fun", theme.Name)
					assert.Equal(t, 7, theme.Month)
					assert.Equal(t, 2025, theme.Year)
				}
			}
		})
	}
}

func TestHandler_DeleteTheme(t *testing.T) {
	themeID := uuid.New()

	tests := []struct {
		name           string
		themeID        string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:    "successful delete theme",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, themeID).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:    "invalid UUID format",
			themeID: "invalid-uuid",
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:    "theme not found",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, themeID).Return(errors.New("theme not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:    "repository error",
			themeID: themeID.String(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, themeID).Return(errors.New("database connection failed"))
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

			req := httptest.NewRequest("DELETE", "/themes/"+tt.themeID, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				// Success case - should have success message
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var successResp map[string]interface{}
				err = json.Unmarshal(body, &successResp)
				assert.NoError(t, err)
				assert.Equal(t, "Theme deleted successfully", successResp["message"])
			}
		})
	}
}
