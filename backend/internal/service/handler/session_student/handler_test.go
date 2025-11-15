package sessionstudent_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/models"
	sessionstudent "specialstandard/internal/service/handler/session_student"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_CreateSessionStudent(t *testing.T) {
	sessionID := uuid.New()
	sessionID2 := uuid.New()
	studentID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockSessionStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful_create_session_student",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": true,
				"notes": "Student participated well in group activities"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(&[]models.SessionStudent{
					{
						SessionID: sessionID,
						StudentID: studentID,
						Present:   true,
						Notes:     stringPtr("Student participated well in group activities"),
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "successful_create_session_student_minimal_data",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": false
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(&[]models.SessionStudent{
					{
						SessionID: sessionID,
						StudentID: studentID,
						Present:   false,
						Notes:     nil,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name:        "invalid_JSON_body",
			requestBody: `{"session_ids": ["` + sessionID.String() + `"], "student_ids": /* missing comma */}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_session_id",
			requestBody: `{"student_ids": ["` + studentID.String() + `"], "present": true}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_student_id",
			requestBody: `{"session_ids": ["` + sessionID.String() + `"], "present": true}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_session_id_format",
			requestBody: `{"student_ids": ["` + studentID.String() + `"], "present": true, "session_ids": "not-a-uuid"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing should fail on invalid UUID format
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_student_id_format",
			requestBody: `{"session_ids": ["` + sessionID.String() + `"], "present": true, "student_ids": "not-a-uuid"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing should fail on invalid UUID format
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "empty_session_id_nil_uuid",
			requestBody: `{
				"session_ids": ["00000000-0000-0000-0000-000000000000"],
				"student_ids": ["42a36e4a-1a3e-4a08-aac0-a1ca769e79d1"],
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// The handler validates zero UUIDs before calling repository, so no mock needed
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "empty_student_id_nil_uuid",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["00000000-0000-0000-0000-000000000000"],
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// The handler validates zero UUIDs before calling repository, so no mock needed
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "duplicate_session_student_relationship",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": true,
				"notes": "Duplicate entry"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(nil, errors.New("duplicate key value violates unique constraint"))
			},
			expectedStatus: fiber.StatusConflict,
			wantErr:        true,
		},
		{
			name: "repository_save_error",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:        "empty_JSON_body",
			requestBody: ``,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No repo call expected
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "unique_violation_conflict_error",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(nil, errors.New("pq: unique_violation: duplicate student in session"))
			},
			expectedStatus: fiber.StatusConflict,
			wantErr:        true,
		},
		{
			name: "successful_multiple_session_ids",
			requestBody: `{
				"session_ids": ["` + sessionID.String() + `", "` + sessionID2.String() + `"],
				"student_ids": ["` + studentID.String() + `"],
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetDB").Return((*pgxpool.Pool)(nil))
				m.On("CreateSessionStudent",
					mock.AnythingOfType("*fasthttp.RequestCtx"),
					(*pgxpool.Pool)(nil),
					mock.AnythingOfType("*models.CreateSessionStudentInput"),
				).Return(&[]models.SessionStudent{
					{
						SessionID: sessionID,
						StudentID: studentID,
						Present:   true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						SessionID: sessionID2,
						StudentID: studentID,
						Present:   true,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockRepo)

			handler := sessionstudent.NewHandler(mockRepo)
			app.Post("/session_students", handler.CreateSessionStudent)

			req := httptest.NewRequest("POST", "/session_students", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusCreated {
				// Success case - validate created session student data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var sessionStudents []models.SessionStudent

				err = json.Unmarshal(body, &sessionStudents)
				assert.NoError(t, err)

				var allSessionIDs []uuid.UUID
				for _, sessionStudent := range sessionStudents {
					allSessionIDs = append(allSessionIDs, sessionStudent.SessionID)
				}

				for _, sessionStudent := range sessionStudents {
					// Validate response data
					assert.Contains(t, allSessionIDs, sessionStudent.SessionID)
					assert.Equal(t, studentID, sessionStudent.StudentID)
					assert.False(t, sessionStudent.CreatedAt.IsZero())
					assert.False(t, sessionStudent.UpdatedAt.IsZero())

					switch tt.name {
					case "successful_create_session_student":
						assert.True(t, sessionStudent.Present)
						assert.NotNil(t, sessionStudent.Notes)
						assert.Contains(t, *sessionStudent.Notes, "participated well")
					case "successful_create_session_student_minimal_data":
						assert.False(t, sessionStudent.Present)
						assert.Nil(t, sessionStudent.Notes)
					}
				}
			}
		})
	}
}

func TestHandler_PatchSessionStudent(t *testing.T) {
	sessionID := uuid.New()
	studentID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockSessionStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful_patch_present_only",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"present": false
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionStudent := &models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   false,
					Notes:     stringPtr("Original notes"),
				}
				ratings := []models.SessionRating{}
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(sessionStudent, ratings, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful_patch_notes_only",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"notes": "Updated notes about student progress"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionStudent := &models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   true,
					Notes:     stringPtr("Updated notes about student progress"),
				}
				ratings := []models.SessionRating{}
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(sessionStudent, ratings, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful_patch_both_fields",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"present": true,
				"notes": "Student showed improvement"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionStudent := &models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   true,
					Notes:     stringPtr("Student showed improvement"),
				}
				ratings := []models.SessionRating{}
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(sessionStudent, ratings, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful_ratings_update_only",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"ratings": [
					{
						"category": "visual_cue",
						"level": "minimal",
						"description": "Student occasionally makes eye contact."
					},
					{
						"category": "verbal_cue",
						"level": "moderate",
						"description": "Student responds to questions with complete sentences."
					}
				]
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				sessionStudent := &models.SessionStudent{
					SessionID: sessionID,
					StudentID: studentID,
					Present:   true,
					Notes:     stringPtr("Original notes"),
				}
				ratings := []models.SessionRating{
					{
						Category:    stringPtr("visual_cue"),
						Level:       stringPtr("minimal"),
						Description: stringPtr("Student occasionally makes eye contact."),
					},
					{
						Category:    stringPtr("verbal_cue"),
						Level:       stringPtr("moderate"),
						Description: stringPtr("Student responds to questions with complete sentences."),
					},
				}
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(sessionStudent, ratings, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "invalid_ratings_category",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"ratings": [
					{
						"category": "invalid_category",
						"level": "minimal",
						"description": "Student occasionally makes eye contact."
					}
				]
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid_ratings_level",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"ratings": [
					{
						"category": "visual_cue",
						"level": "invalid_level",
						"description": "Student occasionally makes eye contact."
					}
				]
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_JSON_body",
			requestBody: `{"session_id": "` + sessionID.String() + `", "present": /* missing value */}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_session_id",
			requestBody: `{"student_id": "` + studentID.String() + `", "present": true}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_student_id",
			requestBody: `{"session_id": "` + sessionID.String() + `", "present": true}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "session_student_not_found",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(nil, nil, errors.New("no rows affected"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "foreign_key_violation",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"present": true
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(nil, nil, errors.New("foreign key violation"))
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "repository_error",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `",
				"present": false
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("RateStudentSession", mock.Anything, mock.AnythingOfType("*models.PatchSessionStudentInput")).Return(nil, nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(mocks.MockSessionStudentRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := sessionstudent.NewHandler(mockRepo)
			app.Patch("/session_students", handler.PatchStudentSessionRatings)

			req := httptest.NewRequest("PATCH", "/session_students", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.mockSetup != nil {
				mockRepo.AssertExpectations(t)
			}

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				assert.Equal(t, sessionID.String(), result["sessionId"])
				assert.Equal(t, studentID.String(), result["studentId"])

				switch tt.name {
				case "successful_patch_present_only":
					assert.False(t, result["present"].(bool))
					assert.NotNil(t, result["notes"])
				case "successful_patch_notes_only":
					assert.True(t, result["present"].(bool))
					assert.NotNil(t, result["notes"])
					assert.Contains(t, result["notes"].(string), "Updated notes")
				case "successful_ratings_update_only":
					assert.True(t, result["present"].(bool))
					ratings := result["ratings"].([]interface{})
					assert.Len(t, ratings, 2)
				case "successful_patch_both_fields":
					assert.True(t, result["present"].(bool))
					assert.NotNil(t, result["notes"])
					assert.Contains(t, result["notes"].(string), "improvement")
				}
			}

			if tt.wantErr {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				if err == nil {
					assert.Contains(t, errorResp, "error")
				}
			}
		})
	}
}

func TestHandler_DeleteSessionStudent(t *testing.T) {
	sessionID := uuid.New()
	studentID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockSessionStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful_delete_session_student",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(nil)
			},
			expectedStatus: fiber.StatusNoContent,
			wantErr:        false,
		},
		{
			name:        "invalid_JSON_body",
			requestBody: `{"session_id": "` + sessionID.String() + `", "student_id": /* missing comma */}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_session_id",
			requestBody: `{"student_id": "` + studentID.String() + `"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing_student_id",
			requestBody: `{"session_id": "` + sessionID.String() + `"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_session_id_format",
			requestBody: `{"student_id": "` + studentID.String() + `", "session_id": "not-a-uuid"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing should succeed but UUID validation should fail
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_student_id_format",
			requestBody: `{"session_id": "` + sessionID.String() + `", "student_id": "not-a-uuid"}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - JSON parsing should succeed but UUID validation should fail
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "empty_session_id_nil_uuid",
			requestBody: `{
				"session_id": "00000000-0000-0000-0000-000000000000",
				"student_id": "` + studentID.String() + `"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// The nil UUID should be caught by validation before repository call
				// But add .Maybe() in case validation logic changes
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(errors.New("session not found")).Maybe()
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "empty_student_id_nil_uuid",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "00000000-0000-0000-0000-000000000000"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// The nil UUID should be caught by validation before repository call
				// But add .Maybe() in case validation logic changes
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(errors.New("student not found")).Maybe()
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "session_student_relationship_not_found",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(errors.New("no rows affected"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name: "repository_delete_error",
			requestBody: `{
				"session_id": "` + sessionID.String() + `",
				"student_id": "` + studentID.String() + `"
			}`,
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("DeleteSessionStudent", mock.Anything, mock.AnythingOfType("*models.DeleteSessionStudentInput")).Return(errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(mocks.MockSessionStudentRepository)
			tt.mockSetup(mockRepo)

			handler := sessionstudent.NewHandler(mockRepo)
			app.Delete("/session_students", handler.DeleteSessionStudent)

			req := httptest.NewRequest("DELETE", "/session_students", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			if !tt.wantErr && resp.StatusCode == fiber.StatusNoContent {
				// Success case - should have empty response body
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Empty(t, body)
			}

			if tt.wantErr {
				// Error cases - validate error response structure
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp, "error")

				// Validate specific error messages
				switch tt.name {
				case "invalid_JSON_body":
					assert.Contains(t, errorResp["error"], "Invalid JSON format")
				case "missing_session_id", "empty_session_id_nil_uuid":
					assert.Contains(t, errorResp["error"], "Session ID is required")
				case "missing_student_id", "empty_student_id_nil_uuid":
					assert.Contains(t, errorResp["error"], "Student ID is required")
				case "invalid_session_id_format":
					// This might be caught by JSON parsing or UUID validation
					errorMsg := errorResp["error"].(string)
					assert.True(t, strings.Contains(errorMsg, "Invalid JSON format") || strings.Contains(errorMsg, "Invalid") || strings.Contains(errorMsg, "Session ID"))
				case "invalid_student_id_format":
					// This might be caught by JSON parsing or UUID validation
					errorMsg := errorResp["error"].(string)
					assert.True(t, strings.Contains(errorMsg, "Invalid JSON format") || strings.Contains(errorMsg, "Invalid") || strings.Contains(errorMsg, "Student ID"))
				case "session_student_relationship_not_found":
					assert.Contains(t, errorResp["error"], "Session student relationship not found")
				case "repository_delete_error":
					assert.Contains(t, errorResp["error"], "Failed to delete session student")
				}
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func TestHandler_GetStudentAttendance(t *testing.T) {
	studentID := uuid.New()

	tests := []struct {
		name           string
		studentID      string
		queryParams    string
		mockSetup      func(*mocks.MockSessionStudentRepository)
		expectedStatus int
		expectedBody   map[string]interface{}
		wantErr        bool
	}{
		{
			name:        "successful_get_all_attendance",
			studentID:   studentID.String(),
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 18
				totalCount := 20
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(18),
				"total_count":   float64(20),
			},
			wantErr: false,
		},
		{
			name:        "successful_with_date_range",
			studentID:   studentID.String(),
			queryParams: "?date_from=2024-01-01&date_to=2024-12-31",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 10
				totalCount := 12
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(10),
				"total_count":   float64(12),
			},
			wantErr: false,
		},
		{
			name:        "successful_with_only_date_from",
			studentID:   studentID.String(),
			queryParams: "?date_from=2024-01-01",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 15
				totalCount := 18
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(15),
				"total_count":   float64(18),
			},
			wantErr: false,
		},
		{
			name:        "successful_with_only_date_to",
			studentID:   studentID.String(),
			queryParams: "?date_to=2024-12-31",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 8
				totalCount := 10
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(8),
				"total_count":   float64(10),
			},
			wantErr: false,
		},
		{
			name:        "student_with_no_sessions",
			studentID:   studentID.String(),
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 0
				totalCount := 0
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(0),
				"total_count":   float64(0),
			},
			wantErr: false,
		},
		{
			name:        "invalid_student_id_format",
			studentID:   "not-a-uuid",
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - UUID validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "nil_uuid_student_id",
			studentID:   "00000000-0000-0000-0000-000000000000",
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// The nil UUID should be caught by validation
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_date_from_format",
			studentID:   studentID.String(),
			queryParams: "?date_from=invalid-date",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - date parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid_date_to_format",
			studentID:   studentID.String(),
			queryParams: "?date_to=2024-13-45", // Invalid month and day
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				// No mock setup needed - date parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "repository_error",
			studentID:   studentID.String(),
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(nil, nil, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:        "student_not_found",
			studentID:   uuid.New().String(), // Different UUID
			queryParams: "",
			mockSetup: func(m *mocks.MockSessionStudentRepository) {
				presentCount := 0
				totalCount := 0
				m.On("GetStudentAttendance",
					mock.Anything,
					mock.Anything,
				).Return(&presentCount, &totalCount, nil)
			},
			expectedStatus: fiber.StatusOK,
			expectedBody: map[string]interface{}{
				"present_count": float64(0),
				"total_count":   float64(0),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(mocks.MockSessionStudentRepository)

			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			handler := sessionstudent.NewHandler(mockRepo)
			// NOW use :id to match the route exactly
			app.Get("/students/:id/attendance", handler.GetStudentAttendance)

			req := httptest.NewRequest("GET", "/students/"+tt.studentID+"/attendance"+tt.queryParams, nil)
			resp, _ := app.Test(req, -1)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				var result map[string]interface{}
				err = json.Unmarshal(body, &result)
				assert.NoError(t, err)

				if tt.expectedBody != nil {
					for key, expectedValue := range tt.expectedBody {
						actualValue, exists := result[key]
						if key == "student_id" && !exists {
							// student_id might not be in response based on implementation
							continue
						}
						assert.True(t, exists, "Expected key %s to exist", key)
						assert.Equal(t, expectedValue, actualValue, "Value mismatch for key %s", key)
					}
				}
			} else if tt.wantErr && resp.StatusCode >= 400 && resp.StatusCode < 600 {
				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				if err == nil && errorResp["error"] != nil {
					switch tt.name {
					case "invalid_student_id_format":
						assert.Contains(t, errorResp["error"], "invalid")
					case "nil_uuid_student_id":
						assert.Contains(t, errorResp["error"], "invalid")
					case "invalid_date_from_format", "invalid_date_to_format":
						assert.Contains(t, errorResp["error"], "Invalid date")
					case "repository_error":
						assert.Contains(t, errorResp["error"], "Failed to get attendance")
					}
				}
			}

			if tt.mockSetup != nil {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
