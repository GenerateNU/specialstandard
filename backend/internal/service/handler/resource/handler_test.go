package resource_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/resource"
	"specialstandard/internal/storage/mocks"
	"specialstandard/internal/utils"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string     { return &s }
func ptrInt(i int) *int              { return &i }
func ptrTime(t time.Time) *time.Time { return &t }

func TestHandler_PostResource(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name:        "successful post resource",
			requestBody: `{"theme_id": "` + uuid.New().String() + `", "title": "Resource1", "type": "doc"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				createdResource := &models.Resource{
					ID:    uuid.New(),
					Title: ptrString("Resource1"),
					Type:  ptrString("doc"),
				}
				m.On("CreateResource", mock.Anything, mock.Anything).Return(createdResource, nil)
			},
			expectedStatus: fiber.StatusCreated,
		},
		{
			name:        "repository error",
			requestBody: `{"theme_id": "` + uuid.New().String() + `", "title": "Resource1", "type": "doc"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("CreateResource", mock.Anything, mock.Anything).Return((*models.Resource)(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:        "foreign key constraint error - invalid theme",
			requestBody: `{"theme_id": "` + uuid.New().String() + `", "title": "Resource1"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("CreateResource", mock.Anything, mock.Anything).
					Return((*models.Resource)(nil), errs.InvalidRequestData(map[string]string{"theme_id": "invalid theme"}))
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			requestBody:    `{"invalid json`,
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:        "complete resource with all fields",
			requestBody: `{"theme_id": "` + uuid.New().String() + `", "grade_level": 5, "date": "2024-01-15T00:00:00Z", "type": "worksheet", "title": "Math Worksheet", "category": "mathematics", "content": "Addition problems"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("CreateResource", mock.Anything, mock.AnythingOfType("models.ResourceBody")).
					Return(&models.Resource{
						ID:         uuid.New(),
						Title:      ptrString("Math Worksheet"),
						Type:       ptrString("worksheet"),
						GradeLevel: ptrInt(5),
						Category:   ptrString("mathematics"),
						Content:    ptrString("Addition problems"),
					}, nil)
			},
			expectedStatus: fiber.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Post("/resources", handler.PostResource)

			req := httptest.NewRequest("POST", "/resources", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetResource(t *testing.T) {
	resourceID := uuid.New()
	tests := []struct {
		name           string
		resourceID     string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:       "successful_get_resource",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResourceByID", mock.Anything, resourceID).Return(&models.ResourceWithTheme{
					Resource: models.Resource{
						ID:    resourceID,
						Title: ptrString("Resource1"),
						Type:  ptrString("doc"),
					},
					Theme: models.ThemeInfo{
						Name:      "Theme1",
						Month:     6,
						Year:      2025,
						CreatedAt: nil,
						UpdatedAt: nil,
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:       "resource not found - pgx.ErrNoRows",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResourceByID", mock.Anything, resourceID).Return((*models.ResourceWithTheme)(nil), pgx.ErrNoRows)
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:       "repository_error",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResourceByID", mock.Anything, resourceID).Return((*models.ResourceWithTheme)(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "invalid UUID format",
			resourceID:     "not-a-uuid",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Get("/resources/:id", handler.GetResourceByID)

			req := httptest.NewRequest("GET", "/resources/"+tt.resourceID, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				var res models.ResourceWithTheme
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)
				assert.Equal(t, resourceID, res.ID)
			}
		})
	}
}

func TestHandler_GetResources(t *testing.T) {
	themeID := uuid.New()
	testDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name: "successful_get_resources with default pagination",
			url:  "",
			mockSetup: func(m *mocks.MockResourceRepository) {
				resources := []models.ResourceWithTheme{
					{Resource: models.Resource{ID: uuid.New(), Title: ptrString("Resource1"), Type: ptrString("doc")},
						Theme: models.ThemeInfo{Name: "Spring", Month: 3, Year: 2024, CreatedAt: nil, UpdatedAt: nil}},
				}

				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return(resources, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by theme_id",
			url:  "?theme_id=" + themeID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, themeID, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by grade_level",
			url:  "?grade_level=5",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "5", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by type",
			url:  "?type=worksheet",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "worksheet", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by title",
			url:  "?title=Math",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "Math", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by category",
			url:  "?category=mathematics",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "mathematics", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by content",
			url:  "?content=fractions",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "fractions", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by date",
			url:  "?date=" + testDate.Format(time.RFC3339),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", mock.AnythingOfType("*time.Time"), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "filter by theme_name (ILIKE search)",
			url:  "?theme_name=Spring",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "Spring", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "theme_params_filter",
			url:  "?theme_name=Spring&theme_month=3&theme_year=2024",
			mockSetup: func(m *mocks.MockResourceRepository) {
				themeMonth := 3
				themeYear := 2024
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "Spring", (*time.Time)(nil), &themeMonth, &themeYear, utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "multiple filters combined",
			url:  "?theme_id=" + themeID.String() + "&grade_level=5&type=worksheet&category=mathematics",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, themeID, "5", "worksheet", "", "mathematics", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "empty_resources_list",
			url:  "",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "invalid theme_id UUID",
			url:            "?theme_id=invalid-uuid",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid theme month",
			url:            "?theme_month=13",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid theme month - zero",
			url:            "?theme_month=0",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid theme year - too low",
			url:            "?theme_year=1999",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid theme year - too high",
			url:            "?theme_year=2501",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid date format",
			url:            "?date=not-a-date",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "repository_error",
			url:  "",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), utils.NewPagination()).Return([]models.ResourceWithTheme(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=xyz",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name: "Pagination Parameters",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockResourceRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), pagination).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "Large pagination limit",
			url:  "?limit=100",
			mockSetup: func(m *mocks.MockResourceRepository) {
				pagination := utils.Pagination{
					Page:  1,
					Limit: 100,
				}
				m.On("GetResources", mock.Anything, uuid.Nil, "", "", "", "", "", "", (*time.Time)(nil), (*int)(nil), (*int)(nil), pagination).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "invalid grade level - negative",
			url:            "?grade_level=-1",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid grade level - too high",
			url:            "?grade_level=13",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Get("/resources", handler.GetResources)

			req := httptest.NewRequest("GET", "/resources"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchResource(t *testing.T) {
	resourceID := uuid.New()
	newThemeID := uuid.New()
	updateDate := time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		resourceID     string
		requestBody    string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name:        "successful update - title only",
			resourceID:  resourceID.String(),
			requestBody: `{"title": "Updated Title"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.MatchedBy(func(body models.UpdateResourceBody) bool {
					return body.Title != nil && *body.Title == "Updated Title"
				})).Return(&models.Resource{
					ID:    resourceID,
					Title: ptrString("Updated Title"),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "successful update - all fields",
			resourceID:  resourceID.String(),
			requestBody: `{"theme_id": "` + newThemeID.String() + `", "grade_level": 6, "date": "` + updateDate.Format(time.RFC3339) + `", "type": "video", "title": "Updated Title", "category": "science", "content": "New content"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.AnythingOfType("models.UpdateResourceBody")).
					Return(&models.Resource{
						ID:         resourceID,
						ThemeID:    newThemeID,
						GradeLevel: ptrInt(6),
						Date:       ptrTime(updateDate),
						Type:       ptrString("video"),
						Title:      ptrString("Updated Title"),
						Category:   ptrString("science"),
						Content:    ptrString("New content"),
					}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:        "no fields to update error",
			resourceID:  resourceID.String(),
			requestBody: `{}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.AnythingOfType("models.UpdateResourceBody")).
					Return((*models.Resource)(nil), errors.New("no fields to update"))
			},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid UUID",
			resourceID:     "not-a-uuid",
			requestBody:    `{"title": "Updated"}`,
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			resourceID:     resourceID.String(),
			requestBody:    `{"invalid json`,
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:        "repository_error",
			resourceID:  resourceID.String(),
			requestBody: `{"title": "Updated"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.Anything).Return((*models.Resource)(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:        "update grade level only",
			resourceID:  resourceID.String(),
			requestBody: `{"grade_level": 8}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.MatchedBy(func(body models.UpdateResourceBody) bool {
					return body.GradeLevel != nil && *body.GradeLevel == 8
				})).Return(&models.Resource{
					ID:         resourceID,
					GradeLevel: ptrInt(8),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "update with invalid theme_id format",
			resourceID:     resourceID.String(),
			requestBody:    `{"theme_id": "invalid-uuid"}`,
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:        "partial update - category and content",
			resourceID:  resourceID.String(),
			requestBody: `{"category": "updated-category", "content": "updated content"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, resourceID, mock.MatchedBy(func(body models.UpdateResourceBody) bool {
					return body.Category != nil && *body.Category == "updated-category" &&
						body.Content != nil && *body.Content == "updated content"
				})).Return(&models.Resource{
					ID:       resourceID,
					Category: ptrString("updated-category"),
					Content:  ptrString("updated content"),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Patch("/resources/:id", handler.UpdateResource)

			req := httptest.NewRequest("PATCH", "/resources/"+tt.resourceID, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteResource(t *testing.T) {
	resourceID := uuid.New()
	tests := []struct {
		name           string
		resourceID     string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name:       "successful_delete_resource",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("DeleteResource", mock.Anything, resourceID).Return(nil)
			},
			expectedStatus: fiber.StatusNoContent,
		},
		{
			name:       "repository_error",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("DeleteResource", mock.Anything, resourceID).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
		{
			name:           "invalid UUID format",
			resourceID:     "invalid-uuid",
			mockSetup:      func(m *mocks.MockResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
		},
		{
			name:       "foreign key constraint error",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("DeleteResource", mock.Anything, resourceID).
					Return(errs.InternalServerError("foreign key constraint violation"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Delete("/resources/:id", handler.DeleteResource)

			req := httptest.NewRequest("DELETE", "/resources/"+tt.resourceID, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

// Additional test for complex filtering scenarios
func TestHandler_GetResources_ComplexQueries(t *testing.T) {
	themeID := uuid.New()
	testDate := time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)
	themeMonth := 3
	themeYear := 2024

	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name: "all filters combined with pagination",
			url:  "?theme_id=" + themeID.String() + "&grade_level=5&type=worksheet&title=Math&category=mathematics&content=fractions&theme_name=Spring&date=" + testDate.Format(time.RFC3339) + "&theme_month=3&theme_year=2024&page=2&limit=20",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources",
					mock.Anything,
					themeID,
					"5",
					"worksheet",
					"Math",
					"mathematics",
					"fractions",
					"Spring",
					mock.AnythingOfType("*time.Time"),
					&themeMonth,
					&themeYear,
					utils.Pagination{Page: 2, Limit: 20},
				).Return([]models.ResourceWithTheme{
					{
						Resource: models.Resource{
							ID:         uuid.New(),
							ThemeID:    themeID,
							GradeLevel: ptrInt(5),
							Date:       ptrTime(testDate),
							Type:       ptrString("worksheet"),
							Title:      ptrString("Math Fractions"),
							Category:   ptrString("mathematics"),
							Content:    ptrString("fractions exercises"),
						},
						Theme: models.ThemeInfo{
							Name:  "Spring",
							Month: 3,
							Year:  2024,
						},
					},
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "empty result with valid filters",
			url:  "?category=non-existent-category",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources",
					mock.Anything,
					uuid.Nil,
					"",
					"",
					"",
					"non-existent-category",
					"",
					"",
					(*time.Time)(nil),
					(*int)(nil),
					(*int)(nil),
					utils.NewPagination(),
				).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockResourceRepository)
			tt.mockSetup(mockRepo)

			handler := resource.NewHandler(mockRepo)
			app.Get("/resources", handler.GetResources)

			req := httptest.NewRequest("GET", "/resources"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				var resources []models.ResourceWithTheme
				err = json.Unmarshal(body, &resources)
				assert.NoError(t, err)
			}
		})
	}
}
