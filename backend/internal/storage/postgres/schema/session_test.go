package schema_test

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/service/handler/session"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	// Generate a test therapist ID to use across tests
	testTherapistID := uuid.New()
	
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockSessionRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get sessions with default pagination",
			url:  fmt.Sprintf("?therapist_id=%s", testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository) {
				sessions := []models.Session{
					{
						ID:            uuid.New(),
						TherapistID:   testTherapistID,
						StartDateTime: time.Now(),
						EndDateTime:   time.Now().Add(time.Hour),
						Notes:         ptrString("Test session"),
						CreatedAt:     ptrTime(time.Now()),
						UpdatedAt:     ptrTime(time.Now()),
					},
				}
				// Updated to include filter (nil) and therapistID parameters
				m.On("GetSessions", mock.Anything, utils.NewPagination(), (*models.GetSessionRepositoryRequest)(nil), testTherapistID).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  fmt.Sprintf("?therapist_id=%s", testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository) {
				m.On("GetSessions", mock.Anything, utils.NewPagination(), (*models.GetSessionRepositoryRequest)(nil), testTherapistID).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Missing therapist ID",
			url:            "?page=1&limit=10",
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid therapist ID format",
			url:            "?therapist_id=invalid-uuid",
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            fmt.Sprintf("?therapist_id=%s&page=0&limit=-1", testTherapistID.String()),
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            fmt.Sprintf("?therapist_id=%s&page=abc&limit=-1", testTherapistID.String()),
			mockSetup:      func(m *mocks.MockSessionRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name: "Pagination Parameters",
			url:  fmt.Sprintf("?therapist_id=%s&page=2&limit=5", testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository) {
				pagination := utils.Pagination{
					Page:  2,
					Limit: 5,
				}
				m.On("GetSessions", mock.Anything, pagination, (*models.GetSessionRepositoryRequest)(nil), testTherapistID).Return([]models.Session{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "With date filters",
			url:  fmt.Sprintf("?therapist_id=%s&month=11&year=2024", testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository) {
				month := 11
				year := 2024
				filter := &models.GetSessionRepositoryRequest{
					Month: &month,
					Year:  &year,
				}
				m.On("GetSessions", mock.Anything, utils.NewPagination(), filter, testTherapistID).Return([]models.Session{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "With student IDs filter",
			url:  fmt.Sprintf("?therapist_id=%s&student_ids[]=%s&student_ids[]=%s", testTherapistID.String(), uuid.New().String(), uuid.New().String()),
			mockSetup: func(m *mocks.MockSessionRepository) {
				// For student IDs, we expect the filter to have StudentIDs populated
				m.On("GetSessions", mock.Anything, utils.NewPagination(), mock.MatchedBy(func(filter *models.GetSessionRepositoryRequest) bool {
					return filter != nil && filter.StudentIDs != nil && len(*filter.StudentIDs) == 2
				}), testTherapistID).Return([]models.Session{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo)

			mockRepoSSR := new(mocks.MockSessionStudentRepository)
			handler := session.NewHandler(mockRepo, mockRepoSSR)
			app.Get("/sessions", handler.GetSessions)

			req := httptest.NewRequest("GET", "/sessions"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_DeleteSessions(t *testing.T) {
	tests := []struct {
		id             uuid.UUID
		name           string
		mockSetup      func(*mocks.MockSessionRepository, uuid.UUID)
		expectedStatus int
		wantErr        bool
	}{
		{
			id:   uuid.New(),
			name: "Successful Delete Session",
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSession", mock.Anything, id).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			id:   uuid.New(),
			name: "internal server error",
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				m.On("DeleteSession", mock.Anything, id).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	t.Run("Bad UUID Request", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionRepository)
		mockRepoSSR := new(mocks.MockSessionStudentRepository)

		handler := session.NewHandler(mockRepo, mockRepoSSR)
		app.Delete("/sessions/:id", handler.DeleteSessions)

		req := httptest.NewRequest("DELETE", "/sessions/1234", nil)
		res, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo, tt.id)
			mockRepoSSR := new(mocks.MockSessionStudentRepository)

			handler := session.NewHandler(mockRepo, mockRepoSSR)
			app.Delete("/sessions/:id", handler.DeleteSessions)

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/sessions/%s", tt.id.String()), nil)
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PostSessions(t *testing.T) {
	testTherapistID := uuid.New()
	
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionRepository, *mocks.MockSessionStudentRepository)
		expectedStatusCode int
	}{
		{
			name: "Missing Items, Invalid JSON",
			payload: `{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z"
			}`,
			mockSetup:          func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Empty Values that are Required",
			payload: `{
				"start_datetime": "",
				"end_datetime": "",
				"therapist_id": "00000000-0000-0000-0000-000000000000",
				"notes": ""
			}`,
			mockSetup:          func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			name: "Foreign Key Violation: Therapist ID doesn't exist",
			payload: fmt.Sprintf(`{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z",
				"therapist_id": "%s",
				"notes": "Test FK"
			}`, uuid.New().String()),
			mockSetup: func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {
				m.On("PostSession", mock.Anything, mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).
					Return(nil, errors.New("foreign key violation"))
				m.On("GetDB").Return(nil)
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Start time and end time (check constraint violation)",
			payload: fmt.Sprintf(`{
				"start_datetime": "2025-09-14T11:00:00Z",
				"end_datetime": "2025-09-14T10:00:00Z",
				"therapist_id": "%s",
				"notes": "Check violation"
			}`, testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {
				m.On("PostSession", mock.Anything, mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).
					Return(nil, errors.New("check constraint violation"))
				m.On("GetDB").Return(nil)
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
		{
			name: "Success!",
			payload: fmt.Sprintf(`{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z",
				"therapist_id": "%s",
				"notes": "Test Session"
			}`, testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {
				sessions := []models.Session{
					{
						ID:            uuid.New(),
						TherapistID:   testTherapistID,
						StartDateTime: time.Now(),
						EndDateTime:   time.Now().Add(time.Hour),
						Notes:         ptrString("Test Session"),
						CreatedAt:     ptrTime(time.Now()),
						UpdatedAt:     ptrTime(time.Now()),
					},
				}
				m.On("PostSession", mock.Anything, mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).
					Return(&sessions, nil)
				m.On("GetDB").Return(nil)
			},
			expectedStatusCode: fiber.StatusCreated,
		},
		{
			name: "Database Connection Refused",
			payload: fmt.Sprintf(`{
				"start_datetime": "2025-09-14T10:00:00Z",
				"end_datetime": "2025-09-14T11:00:00Z",
				"therapist_id": "%s",
				"notes": "DB connection test"
			}`, testTherapistID.String()),
			mockSetup: func(m *mocks.MockSessionRepository, ms *mocks.MockSessionStudentRepository) {
				m.On("PostSession", mock.Anything, mock.Anything, mock.AnythingOfType("*models.PostSessionInput")).
					Return(nil, errors.New("connection refused"))
				m.On("GetDB").Return(nil)
			},
			expectedStatusCode: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			mockRepoSSR := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockRepo, mockRepoSSR)

			handler := session.NewHandler(mockRepo, mockRepoSSR)
			app.Post("/sessions", handler.PostSessions)

			req := httptest.NewRequest("POST", "/sessions", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PatchSessions(t *testing.T) {
	tests := []struct {
		id                 uuid.UUID
		name               string
		payload            string
		mockSetup          func(*mocks.MockSessionRepository, uuid.UUID)
		expectedStatusCode int
	}{
		{
			id:                 uuid.New(),
			name:               "Parsing PatchInputSession Error",
			payload:            `{"notes": "Missing quote}`,
			mockSetup:          func(m *mocks.MockSessionRepository, id uuid.UUID) {},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "Given ID not found",
			payload: `{"notes": "Trying to update non-existent"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				patch := &models.PatchSessionInput{
					Notes: ptrString("Trying to update non-existent"),
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, pgx.ErrNoRows)
			},
			expectedStatusCode: fiber.StatusNotFound,
		},
		{
			id:      uuid.New(),
			name:    "foreign key violation - therapist doesn't exist",
			payload: `{"therapist_id": "00000000-0000-0000-0000-000000000999"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				therapistID := uuid.MustParse("00000000-0000-0000-0000-000000000999")
				patch := &models.PatchSessionInput{
					TherapistID: &therapistID,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, errors.New("foreign key"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "check constraint violation",
			payload: `{"start_datetime": "2025-09-14T14:00:00Z", "end_datetime": "2025-09-14T12:00:00Z"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T14:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				patch := &models.PatchSessionInput{
					StartTime: &startTime,
					EndTime:   &endTime,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(nil, errors.New("check constraint"))
			},
			expectedStatusCode: fiber.StatusBadRequest,
		},
		{
			id:      uuid.New(),
			name:    "Successfully changed 1 field",
			payload: `{"notes": "The child seeks to be seen more than they wish to be heard"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				notes := ptrString("The child seeks to be seen more than they wish to be heard")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					Notes: notes,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: time.Now(),
					EndDateTime:   time.Now().Add(time.Hour),
					TherapistID:   uuid.New(),
					Notes:         notes,
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			id:      uuid.New(),
			name:    "Successfully changed multiple fields",
			payload: `{"start_datetime": "2025-09-14T12:00:00Z", "end_datetime": "2025-09-14T13:00:00Z"}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T13:00:00Z")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					StartTime: &startTime,
					EndTime:   &endTime,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: startTime,
					EndDateTime:   endTime,
					TherapistID:   uuid.New(),
					Notes:         ptrString("Rescheduled for convenience"),
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
		{
			id:   uuid.New(),
			name: "Successfully changed all patchable fields",
			payload: `{
				"start_datetime": "2025-09-14T12:00:00Z", 
				"end_datetime": "2025-09-14T13:00:00Z", 
				"therapist_id": "28eedfdc-81e1-44e5-a42c-022dc4c3b64d", 
				"notes": "Starting Over"
			}`,
			mockSetup: func(m *mocks.MockSessionRepository, id uuid.UUID) {
				startTime, _ := time.Parse(time.RFC3339, "2025-09-14T12:00:00Z")
				endTime, _ := time.Parse(time.RFC3339, "2025-09-14T13:00:00Z")
				therapistID := uuid.MustParse("28eedfdc-81e1-44e5-a42c-022dc4c3b64d")
				notes := ptrString("Starting Over")
				createdAt := time.Now().Add(-24 * time.Hour)
				now := time.Now()

				patch := &models.PatchSessionInput{
					StartTime:   &startTime,
					EndTime:     &endTime,
					TherapistID: &therapistID,
					Notes:       notes,
				}

				patchedSession := &models.Session{
					ID:            id,
					StartDateTime: startTime,
					EndDateTime:   endTime,
					TherapistID:   therapistID,
					Notes:         notes,
					CreatedAt:     &createdAt,
					UpdatedAt:     &now,
				}
				m.On("PatchSession", mock.Anything, id, patch).Return(patchedSession, nil)
			},
			expectedStatusCode: fiber.StatusOK,
		},
	}

	t.Run("Bad UUID Request", func(t *testing.T) {
		app := fiber.New(fiber.Config{
			ErrorHandler: errs.ErrorHandler,
		})
		mockRepo := new(mocks.MockSessionRepository)
		mockRepoSSR := new(mocks.MockSessionStudentRepository)

		handler := session.NewHandler(mockRepo, mockRepoSSR)
		app.Patch("/sessions/:id", handler.PatchSessions)

		req := httptest.NewRequest("PATCH", "/sessions/0345", nil)
		req.Header.Set("Content-Type", "application/json")
		res, _ := app.Test(req, -1)

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
		mockRepo.AssertExpectations(t)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionRepository)
			tt.mockSetup(mockRepo, tt.id)

			mockRepoSSR := new(mocks.MockSessionStudentRepository)
			handler := session.NewHandler(mockRepo, mockRepoSSR)
			app.Patch("/sessions/:id", handler.PatchSessions)

			req := httptest.NewRequest("PATCH", "/sessions/"+tt.id.String(), strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t) // ?
		})
	}
}