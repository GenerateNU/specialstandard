package theme_test

import (
	"errors"
	"net/http/httptest"
	"specialstandard/internal/utils"
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
		body           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful create theme",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			}`,
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
			name: "invalid JSON",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			`, // Missing closing brace
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "validation error - missing name",
			body: `{
				"month": 3,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "validation error - invalid month",
			body: `{
				"name": "Spring",
				"month": 13,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "validation error - invalid year",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 1800
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "foreign key constraint error",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(nil, errors.New("foreign key constraint violated"))
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "database connection error",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("CreateTheme", mock.Anything, mock.AnythingOfType("*models.CreateThemeInput")).Return(nil, errors.New("connection refused"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "repository error",
			body: `{
				"name": "Spring",
				"month": 3,
				"year": 2024
			}`,
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

			req := httptest.NewRequest("POST", "/themes", strings.NewReader(tt.body))
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
		url            string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get themes",
			url:  "",
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
				m.On("GetThemes", mock.Anything, utils.NewPagination()).Return(themes, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "",
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemes", mock.Anything, utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockThemeRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockThemeRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name: "Pagination Parameters",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockThemeRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetThemes", mock.Anything, pagination).Return([]models.Theme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
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

			req := httptest.NewRequest("GET", "/themes"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetThemeByID(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get theme by id",
			id:   uuid.NewString(),
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
			name: "invalid UUID format",
			id:   "invalid-uuid",
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "theme not found",
			id:   uuid.NewString(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("GetThemeByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil, errs.NotFound("Error querying database for given ID"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "repository error",
			id:   uuid.NewString(),
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

			req := httptest.NewRequest("GET", "/themes/"+tt.id, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchTheme(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		body           string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful patch theme",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring",
				"month": 4
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				theme := &models.Theme{
					ID:        uuid.New(),
					Name:      "Updated Spring",
					Month:     4,
					Year:      2024,
					CreatedAt: ptrTime(time.Now()),
					UpdatedAt: ptrTime(time.Now()),
				}
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(theme, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "invalid UUID format",
			id:   "invalid-uuid",
			body: `{
				"name": "Updated Spring"
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid JSON",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring",
				"month": 4
			`, // Missing closing brace
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "validation error - invalid month",
			id:   uuid.NewString(),
			body: `{
				"month": 13
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "validation error - invalid year",
			id:   uuid.NewString(),
			body: `{
				"year": 1800
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "no fields provided to update",
			id:   uuid.NewString(),
			body: `{}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errs.BadRequest("No fields provided to update"))
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "theme not found",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring"
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errs.NotFound("error querying database for given theme ID"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "foreign key constraint error",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring"
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("foreign key constraint violated"))
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "database connection error",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring"
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("connection refused"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "repository error",
			id:   uuid.NewString(),
			body: `{
				"name": "Updated Spring"
			}`,
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("PatchTheme", mock.Anything, mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("*models.UpdateThemeInput")).Return(nil, errors.New("database error"))
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

			req := httptest.NewRequest("PATCH", "/themes/"+tt.id, strings.NewReader(tt.body))
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
		id             string
		mockSetup      func(*mocks.MockThemeRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful delete theme",
			id:   uuid.NewString(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "invalid UUID format",
			id:   "invalid-uuid",
			mockSetup: func(m *mocks.MockThemeRepository) {
				// No mock calls expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "theme not found (idempotent)",
			id:   uuid.NewString(),
			mockSetup: func(m *mocks.MockThemeRepository) {
				m.On("DeleteTheme", mock.Anything, mock.AnythingOfType("uuid.UUID")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			id:   uuid.NewString(),
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

			req := httptest.NewRequest("DELETE", "/themes/"+tt.id, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
