package therapist_test

import (
	"errors"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"strings"
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

func TestHandler_GetTherapistByID(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get therapist by id with district",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				districtID := 1
				schoolNames := []string{"Generate Elementary"}
				districtName := "Generate School District"
				therapist := &models.Therapist{
					ID:           uuid.New(),
					FirstName:    "Kevin",
					LastName:     "Matula",
					Email:        "matulakevin91@gmail.com",
					Schools:      []int{1, 2},
					DistrictID:   &districtID,
					SchoolNames:  &schoolNames,
					DistrictName: &districtName,
					Active:       true,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				}
				m.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get therapist by id without district",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "matulakevin91@gmail.com",
					Schools:    []int{},
					DistrictID: nil,
					Active:     true,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "not found error",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("GetTherapistByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errs.NotFound("Therapist not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Get("/therapists/:id", handler.GetTherapistByID)

			req := httptest.NewRequest("GET", "/therapists/9dad94d8-6534-4510-90d7-e4e97c175a65", nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_GetTherapists(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get therapists with mixed district assignments",
			url:  "",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				districtID := 1
				therapists := []models.Therapist{
					{
						ID:         uuid.New(),
						FirstName:  "Kevin",
						LastName:   "Matula",
						Email:      "matulakevin91@gmail.com",
						Schools:    []int{1},
						DistrictID: &districtID,
						Active:     true,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					},
					{
						ID:         uuid.New(),
						FirstName:  "Sarah",
						LastName:   "Mitchell",
						Email:      "sarah.mitchell@example.com",
						Schools:    []int{},
						DistrictID: nil,
						Active:     true,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					},
				}
				m.On("GetTherapists", mock.Anything, utils.NewPagination()).Return(therapists, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("GetTherapists", mock.Anything, utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Pagination with parameters",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetTherapists", mock.Anything, pagination).Return([]models.Therapist{}, nil)
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
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Get("/therapists", handler.GetTherapists)

			req := httptest.NewRequest("GET", "/therapists"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_CreateTherapist(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful create therapist with district and schools",
			body: `{
				"id": "123e4567-e89b-12d3-a456-426614174000",
				"first_name": "Kevin",
				"last_name": "Matula",
				"email": "poop123@gmail.com",
				"schools": [1, 2],
				"district_id": 1
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				districtID := 1
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "poop123@gmail.com",
					Schools:    []int{1, 2},
					DistrictID: &districtID,
					Active:     true,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("CreateTherapist", mock.Anything, mock.AnythingOfType("*models.CreateTherapistInput")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "successful create therapist minimal schools",
			body: `{
				"id": "223e4567-e89b-12d3-a456-426614174000",
				"first_name": "Kevin",
				"last_name": "Matula",
				"email": "poop123@gmail.com",
				"schools": [1],
				"district_id": 1
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "poop123@gmail.com",
					Schools:    []int{},
					DistrictID: nil,
					Active:     true,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("CreateTherapist", mock.Anything, mock.AnythingOfType("*models.CreateTherapistInput")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "invalid json body",
			body:           `{invalid json}`,
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "repository error",
			body: `{
				"id": "323e4567-e89b-12d3-a456-426614174000",
				"first_name": "Kevin",
				"last_name": "Matula",
				"email": "poop123@gmail.com",
				"schools": [1],
				"district_id": 1
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("CreateTherapist", mock.Anything, mock.AnythingOfType("*models.CreateTherapistInput")).Return(nil, errors.New("database error"))
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
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Post("/therapists", handler.CreateTherapist)

			req := httptest.NewRequest("POST", "/therapists", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteTherapist(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful delete therapist by id",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("DeleteTherapist", mock.Anything, mock.AnythingOfType("string")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("DeleteTherapist", mock.Anything, mock.AnythingOfType("string")).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "not found error",
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("DeleteTherapist", mock.Anything, mock.AnythingOfType("string")).Return(errs.NotFound("Therapist not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Delete("/therapists/:id", handler.DeleteTherapist)

			req := httptest.NewRequest("DELETE", "/therapists/4a9a4e58-ea6c-496a-915f-3e8214e77112", nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchTherapist(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockSetup      func(*mocks.MockTherapistRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful patch therapist - update all fields",
			body: `{
				"first_name": "Kevin",
				"last_name": "Matula",
				"email": "poop123@gmail.com",
				"active": false,
				"schools": [1, 2, 3],
				"district_id": 2
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				districtID := 2
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "poop123@gmail.com",
					Active:     false,
					Schools:    []int{1, 2, 3},
					DistrictID: &districtID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("PatchTherapist", mock.Anything, "4a9a4e58-ea6c-496a-915f-3e8214e77112", mock.AnythingOfType("*models.UpdateTherapist")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful patch therapist - partial update email only",
			body: `{"email": "newemail@gmail.com"}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				districtID := 1
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "newemail@gmail.com",
					Active:     true,
					Schools:    []int{1},
					DistrictID: &districtID,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("PatchTherapist", mock.Anything, "4a9a4e58-ea6c-496a-915f-3e8214e77112", mock.AnythingOfType("*models.UpdateTherapist")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful patch therapist - clear district",
			body: `{
				"district_id": null
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapist := &models.Therapist{
					ID:         uuid.New(),
					FirstName:  "Kevin",
					LastName:   "Matula",
					Email:      "matulakevin91@gmail.com",
					Active:     true,
					Schools:    []int{},
					DistrictID: nil,
					CreatedAt:  time.Now(),
					UpdatedAt:  time.Now(),
				}
				m.On("PatchTherapist", mock.Anything, "4a9a4e58-ea6c-496a-915f-3e8214e77112", mock.AnythingOfType("*models.UpdateTherapist")).Return(therapist, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:           "invalid json body",
			body:           `{invalid json}`,
			mockSetup:      func(m *mocks.MockTherapistRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "repository error",
			body: `{
				"first_name": "Kevin"
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("PatchTherapist", mock.Anything, "4a9a4e58-ea6c-496a-915f-3e8214e77112", mock.AnythingOfType("*models.UpdateTherapist")).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name: "not found error",
			body: `{
				"email": "test@example.com"
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				m.On("PatchTherapist", mock.Anything, "4a9a4e58-ea6c-496a-915f-3e8214e77112", mock.AnythingOfType("*models.UpdateTherapist")).Return(nil, errs.NotFound("Therapist not found"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			handler := therapist.NewHandler(mockRepo)
			app.Patch("/therapists/:id", handler.PatchTherapist)

			req := httptest.NewRequest("PATCH", "/therapists/4a9a4e58-ea6c-496a-915f-3e8214e77112", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
