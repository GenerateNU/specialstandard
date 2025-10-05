package resource_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"strings"
	"testing"

	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/resource"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string { return &s }

func TestHandler_PostResource(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name:        "successful post resource",
			requestBody: `{"title": "Resource1", "type": "doc"}`,
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
			requestBody: `{"title": "Resource1", "type": "doc"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("CreateResource", mock.Anything, mock.Anything).Return((*models.Resource)(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
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
				m.On("GetResourceByID", mock.Anything, resourceID).Return(&models.Resource{
					ID:    resourceID,
					Title: ptrString("Resource1"),
					Type:  ptrString("doc"),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:       "repository_error",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResourceByID", mock.Anything, resourceID).Return(&models.Resource{}, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
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
				var res models.Resource
				err = json.Unmarshal(body, &res)
				assert.NoError(t, err)
				assert.Equal(t, resourceID, res.ID)
			}
		})
	}
}

func TestHandler_GetResources(t *testing.T) {
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

				m.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, utils.NewPagination()).Return(resources, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "empty_resources_list",
			url:  "",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "theme_params filter",
			url:  "?theme_name=Spring&theme_month=3&theme_year=2024",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, utils.NewPagination()).Return([]models.ResourceWithTheme{}, nil)
			},
			expectedStatus: fiber.StatusOK,
		},
		{
			name: "repository_error",
			url:  "",
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, utils.NewPagination()).Return([]models.ResourceWithTheme(nil), errors.New("database error"))
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
			url:            "?page=abc&limit=-1",
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
				m.On("GetResources", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, pagination).Return([]models.ResourceWithTheme{}, nil)
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
		})
	}
}

func TestHandler_PatchResource(t *testing.T) {
	resourceID := uuid.New()
	tests := []struct {
		name           string
		resourceID     string
		requestBody    string
		mockSetup      func(*mocks.MockResourceRepository)
		expectedStatus int
	}{
		{
			name:        "repository_error",
			resourceID:  resourceID.String(),
			requestBody: `{"title": "Updated"}`,
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("UpdateResource", mock.Anything, mock.Anything, mock.Anything).Return((*models.Resource)(nil), errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
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
				m.On("DeleteResource", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: fiber.StatusNoContent,
		},
		{
			name:       "repository_error",
			resourceID: resourceID.String(),
			mockSetup: func(m *mocks.MockResourceRepository) {
				m.On("DeleteResource", mock.Anything, mock.Anything).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
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
