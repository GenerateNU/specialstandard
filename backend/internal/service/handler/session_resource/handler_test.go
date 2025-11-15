package session_resource_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/session_resource"
	"specialstandard/internal/storage/mocks"
	"specialstandard/internal/utils"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string {
	return &s
}

func ptrInt(i int) *int {
	return &i
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestHandler_GetSessionResources(t *testing.T) {
	tests := []struct {
		name           string
		sessionID      string
		url            string
		mockSetup      func(*mocks.MockSessionResourceRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:      "successful get resources - multiple resources with default pagination",
			sessionID: uuid.New().String(),
			url:       "",
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				resources := []models.Resource{
					{
						ID:         uuid.New(),
						ThemeID:    uuid.New(),
						GradeLevel: ptrInt(5),
						Week:       ptrInt(2),
						Type:       ptrString("worksheet"),
						Title:      ptrString("Math Worksheet"),
						Category:   ptrString("math"),
						Content:    ptrString("Basic arithmetic"),
						CreatedAt:  ptrTime(time.Now()),
						UpdatedAt:  ptrTime(time.Now()),
					},
					{
						ID:         uuid.New(),
						ThemeID:    uuid.New(),
						GradeLevel: ptrInt(5),
						Week:       ptrInt(2),
						Type:       ptrString("activity"),
						Title:      ptrString("Reading Activity"),
						Category:   ptrString("language"),
						Content:    ptrString("Comprehension exercise"),
						CreatedAt:  ptrTime(time.Now()),
						UpdatedAt:  ptrTime(time.Now()),
					},
				}
				m.On("GetResourcesBySessionID", mock.Anything, mock.AnythingOfType("uuid.UUID"), utils.NewPagination()).Return(resources, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "successful get resources - empty array",
			sessionID: uuid.New().String(),
			url:       "",
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				m.On("GetResourcesBySessionID", mock.Anything, mock.AnythingOfType("uuid.UUID"), utils.NewPagination()).Return([]models.Resource{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "repository returns nil - converts to empty array",
			sessionID: uuid.New().String(),
			url:       "",
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				m.On("GetResourcesBySessionID", mock.Anything, mock.AnythingOfType("uuid.UUID"), utils.NewPagination()).Return(([]models.Resource)(nil), nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "repository error",
			sessionID: uuid.New().String(),
			url:       "",
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				m.On("GetResourcesBySessionID", mock.Anything, mock.AnythingOfType("uuid.UUID"), utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			sessionID:      uuid.New().String(),
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockSessionResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			sessionID:      uuid.New().String(),
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockSessionResourceRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name:      "Pagination Parameters",
			sessionID: uuid.New().String(),
			url:       "?page=2&limit=5",
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetResourcesBySessionID", mock.Anything, mock.AnythingOfType("uuid.UUID"), pagination).Return([]models.Resource{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
	}

	t.Run("Invalid UUID format", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionResourceRepository)

		handler := session_resource.NewHandler(mockRepo)
		app.Get("/sessions/:id/resources", handler.GetSessionResources)

		req := httptest.NewRequest("GET", "/sessions/invalid-uuid/resources", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "GetResourcesBySessionID")
	})

	t.Run("Empty ID parameter", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionResourceRepository)

		handler := session_resource.NewHandler(mockRepo)
		app.Get("/sessions/:id/resources", handler.GetSessionResources)

		req := httptest.NewRequest("GET", "/sessions//resources", nil)
		resp, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
		mockRepo.AssertNotCalled(t, "GetResourcesBySessionID")
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionResourceRepository)
			tt.mockSetup(mockRepo)

			handler := session_resource.NewHandler(mockRepo)
			app.Get("/sessions/:id/resources", handler.GetSessionResources)

			req := httptest.NewRequest("GET", fmt.Sprintf("/sessions/%s/resources", tt.sessionID)+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if resp.StatusCode == fiber.StatusOK {
				var resources []models.Resource
				err := json.NewDecoder(resp.Body).Decode(&resources)
				assert.NoError(t, err)
				assert.NotNil(t, resources) // Should never be nil, always array
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PostSessionResource(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionResourceRepository)
		expectedStatusCode int
	}{
		{
			name: "Invalid JSON",
			payload: `{
				"session_id": "missing closing quote
			}`,
			mockSetup:          func(m *mocks.MockSessionResourceRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Invalid UUID format in session_id",
			payload: `{
				"session_id": "not-a-uuid",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup:          func(m *mocks.MockSessionResourceRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Invalid UUID format in resource_id",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "not-a-uuid"
			}`,
			mockSetup:          func(m *mocks.MockSessionResourceRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Foreign key violation - session not found",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.CreateSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("PostSessionResource", mock.Anything, req).Return(nil, errors.New("ERROR: insert or update on table \"session_resource\" violates foreign key constraint (SQLSTATE 23503)"))
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			name: "Check constraint violation",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.CreateSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("PostSessionResource", mock.Anything, req).Return(nil, errors.New("violates check constraint"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Database connection error",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.CreateSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("PostSessionResource", mock.Anything, req).Return(nil, errors.New("connection refused"))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Generic database error",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.CreateSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("PostSessionResource", mock.Anything, req).Return(nil, errors.New("some database error"))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Success!",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")

				req := models.CreateSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}

				response := &models.SessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("PostSessionResource", mock.Anything, req).Return(response, nil)
			},
			expectedStatusCode: fiber.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionResourceRepository)
			tt.mockSetup(mockRepo)

			handler := session_resource.NewHandler(mockRepo)
			app.Post("/session-resource", handler.PostSessionResource)

			req := httptest.NewRequest("POST", "/session-resource", bytes.NewReader([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteSessionResource(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionResourceRepository)
		expectedStatusCode int
	}{
		{
			name: "Invalid JSON",
			payload: `{
				"session_id": "missing closing quote
			}`,
			mockSetup:          func(m *mocks.MockSessionResourceRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Invalid UUID format",
			payload: `{
				"session_id": "not-a-uuid",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup:          func(m *mocks.MockSessionResourceRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Foreign key violation - session or resource not found",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.DeleteSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("DeleteSessionResource", mock.Anything, req).Return(errors.New("ERROR: violates foreign key constraint (SQLSTATE 23503)"))
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			name: "Relationship not found",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.DeleteSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("DeleteSessionResource", mock.Anything, req).Return(errors.New("relationship not found"))
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			name: "Database connection error",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.DeleteSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("DeleteSessionResource", mock.Anything, req).Return(errors.New("connection refused"))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Generic database error",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.DeleteSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("DeleteSessionResource", mock.Anything, req).Return(errors.New("some database error"))
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Success!",
			payload: `{
				"session_id": "123e4567-e89b-12d3-a456-426614174000",
				"resource_id": "123e4567-e89b-12d3-a456-426614174001"
			}`,
			mockSetup: func(m *mocks.MockSessionResourceRepository) {
				sessionID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				resourceID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174001")
				req := models.DeleteSessionResource{
					SessionID:  sessionID,
					ResourceID: resourceID,
				}
				m.On("DeleteSessionResource", mock.Anything, req).Return(nil)
			},
			expectedStatusCode: fiber.StatusNoContent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionResourceRepository)
			tt.mockSetup(mockRepo)

			handler := session_resource.NewHandler(mockRepo)
			app.Delete("/session-resource", handler.DeleteSessionResource)

			req := httptest.NewRequest("DELETE", "/session-resource", bytes.NewReader([]byte(tt.payload)))
			req.Header.Set("Content-Type", "application/json")

			resp, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatusCode, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
