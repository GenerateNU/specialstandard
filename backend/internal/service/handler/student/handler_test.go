package student_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/student"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
)

func ptrString(s string) *string {
	return &s
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

func TestHandler_GetStudents(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get students",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:          uuid.New(),
						FirstName:   "Test",
						LastName:    "Student",
						DOB:         time.Now().AddDate(-10, 0, 0),
						TherapistID: uuid.New(),
						Grade:       "Test Grade",
						IEP:         "Test IEP",
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.On("GetStudents", mock.Anything).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty students list",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Get("/students", handler.GetStudents)

			req := httptest.NewRequest("GET", "/students", nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			// Response body validation (new addition)
			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var students []models.Student
				err = json.Unmarshal(body, &students)
				assert.NoError(t, err)

				if tt.name == "empty students list" {
					assert.Len(t, students, 0)
				} else if tt.name == "successful get students" {
					assert.Len(t, students, 1)
					assert.Equal(t, "Test", students[0].FirstName)
					assert.Equal(t, "Student", students[0].LastName)
					assert.Equal(t, "Test Grade", students[0].Grade)
					assert.Equal(t, "Test IEP", students[0].IEP)
				}
			}
		})
	}
}

func TestHandler_GetStudent(t *testing.T) {
	studentID := uuid.New()

	tests := []struct {
		name           string
		studentID      string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:      "successful get student",
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				student := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Now().AddDate(-10, 0, 0),
					TherapistID: uuid.New(),
					Grade:       "Test Grade",
					IEP:         "Test IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(student, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "student not found",
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudent", mock.Anything, studentID).Return(models.Student{}, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:      "invalid UUID format",
			studentID: "invalid-uuid",
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:      "repository error",
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudent", mock.Anything, studentID).Return(models.Student{}, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Get("/students/:id", handler.GetStudent)

			// Make request
			req := httptest.NewRequest("GET", "/students/"+tt.studentID, nil)
			resp, _ := app.Test(req, -1)

			// Basic assertions
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			//////////////////////////////////////////////////
			// Response body validation
			//////////////////////////////////////////////////
			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				// Success case - validate student data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var student models.Student
				err = json.Unmarshal(body, &student)
				assert.NoError(t, err)

				// Validate the student data
				assert.Equal(t, "Test", student.FirstName)
				assert.Equal(t, "Student", student.LastName)
				assert.Equal(t, "Test Grade", student.Grade)
				assert.Equal(t, "Test IEP", student.IEP)
				assert.Equal(t, studentID, student.ID)
			}

			if tt.wantErr {
				// Error cases - validate error response structure
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp, "error")

				switch tt.name {
				case "invalid UUID format":
					assert.Contains(t, errorResp["error"], "Invalid UUID format")
				case "student not found":
					assert.Contains(t, errorResp["error"], "Student not found")
				case "repository error":
					assert.Contains(t, errorResp["error"], "Database error")
				}
			}
		})
	}
}

func TestHandler_UpdateStudent(t *testing.T) {
	studentID := uuid.New()
	therapistID := uuid.New()

	tests := []struct {
		name           string
		studentID      string
		requestBody    string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		// Happy Path - Single Field Updates
		{
			name:      "update grade only",
			studentID: studentID.String(),
			requestBody: `{
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "update IEP only",
			studentID: studentID.String(),
			requestBody: `{
				"iep": "Updated IEP with math accommodations"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},

		// Happy Path - Multiple Field Updates
		{
			name:      "update name and grade",
			studentID: studentID.String(),
			requestBody: `{
				"first_name": "Updated",
				"last_name": "TestStudent",
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},

		// Happy Path - Complex Field Updates
		{
			name:      "update DOB with valid date",
			studentID: studentID.String(),
			requestBody: `{
				"dob": "2010-05-15"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},

		// Error Cases - URL Parameter Errors
		{
			name:      "empty student ID",
			studentID: "",
			requestBody: `{
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},
		{
			name:      "invalid UUID format",
			studentID: "invalid-uuid",
			requestBody: `{
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},

		// Error Cases - Student Existence
		{
			name:      "student not found",
			studentID: studentID.String(),
			requestBody: `{
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudent", mock.Anything, studentID).Return(models.Student{}, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusNotFound,
			wantErr:        true,
		},

		// Error Cases - Request Body Validation
		{
			name:      "invalid JSON body",
			studentID: studentID.String(),
			requestBody: `{
				"grade": "5th"
				// missing comma - invalid JSON
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},

		// Error Cases - Field Validation  
		{
			name:      "invalid date format",
			studentID: studentID.String(),
			requestBody: `{
				"dob": "invalid-date"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:      "invalid therapist UUID",
			studentID: studentID.String(),
			requestBody: `{
				"therapist_id": "bad-uuid"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},

		// Error Cases - Repository Errors
		{
			name:      "UpdateStudent repository error",
			studentID: studentID.String(),
			requestBody: `{
				"grade": "5th"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC),
					TherapistID: therapistID,
					Grade:       "4th",
					IEP:         "Original IEP",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Patch("/students/:id", handler.UpdateStudent)

			// Make request - handle empty ID case
			url := "/students/" + tt.studentID
			req := httptest.NewRequest("PATCH", url, strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			// Assert
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_AddStudent(t *testing.T) {
	therapistID := uuid.New()

	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		// Happy Path
		{
			name: "successful create student",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"dob": "2010-05-15",
				"therapist_id": "` + therapistID.String() + `",
				"grade": "5th",
				"iep": "Active IEP with speech therapy goals"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name: "successful create student with different data",
			requestBody: `{
				"first_name": "Emma",
				"last_name": "Johnson", 
				"dob": "2012-03-22",
				"therapist_id": "` + therapistID.String() + `",
				"grade": "3rd",
				"iep": "Math accommodations and extended time"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},

		// Error Cases - Request Body Validation
		{
			name:        "invalid JSON body",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe"
				// missing comma - invalid JSON
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "empty JSON body",
			requestBody: `{}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - date parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},

		// Error Cases - Field Validation
		{
			name: "invalid date format",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"dob": "invalid-date",
				"therapist_id": "` + therapistID.String() + `",
				"grade": "5th",
				"iep": "Active IEP"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - date parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid therapist UUID format",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"dob": "2010-05-15",
				"therapist_id": "invalid-uuid",
				"grade": "5th", 
				"iep": "Active IEP"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "missing required fields",
			requestBody: `{
				"first_name": "John"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},

		// Error Cases - Repository Errors
		{
			name: "repository save error",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"dob": "2010-05-15",
				"therapist_id": "` + therapistID.String() + `",
				"grade": "5th",
				"iep": "Active IEP"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},

		// Edge Cases
		{
			name: "valid date edge cases", 
			requestBody: `{
				"first_name": "Test",
				"last_name": "Student",
				"dob": "2000-02-29",
				"therapist_id": "` + therapistID.String() + `",
				"grade": "12th",
				"iep": "Graduation accommodations"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			app := fiber.New()
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Post("/students", handler.AddStudent)

			// Make request
			req := httptest.NewRequest("POST", "/students", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			//////////////////////////////////////////////////
			// Response body validation
			//////////////////////////////////////////////////
			if !tt.wantErr && resp.StatusCode == fiber.StatusCreated {
				// Success case - validate created student data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var student models.Student
				err = json.Unmarshal(body, &student)
				assert.NoError(t, err)

				// student data based on request
				switch tt.name {
				case "successful create student":
					assert.Equal(t, "John", student.FirstName)
					assert.Equal(t, "Doe", student.LastName)
					assert.Equal(t, "5th", student.Grade)
					assert.Contains(t, student.IEP, "speech therapy")
				case "successful create student with different data":
					assert.Equal(t, "Emma", student.FirstName)
					assert.Equal(t, "Johnson", student.LastName)
					assert.Equal(t, "3rd", student.Grade)
					assert.Contains(t, student.IEP, "Math accommodations")
				case "valid date edge cases":
					assert.Equal(t, "Test", student.FirstName)
					assert.Equal(t, "12th", student.Grade)
				}

				// Validate that UUID was generated
				assert.NotEqual(t, uuid.Nil, student.ID)
				assert.Equal(t, therapistID, student.TherapistID)

				// Validate date was parsed correctly
				assert.False(t, student.DOB.IsZero())
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
				case "invalid date format":
					assert.Contains(t, errorResp["error"], "Invalid date format")
				case "invalid therapist UUID format":
					assert.Contains(t, errorResp["error"], "Invalid therapist ID format")
				case "missing required fields":
					assert.Contains(t, errorResp["error"], "Invalid date format")
				}
			}
			//////////////////////////////////////////////////
		})
	}
}

func TestHandler_DeleteStudent(t *testing.T) {
	studentID := uuid.New()

	tests := []struct {
		name           string
		studentID      string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:      "successful delete student",
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("DeleteStudent", mock.Anything, studentID).Return(nil)
			},
			expectedStatus: fiber.StatusNoContent,
			wantErr:        false,
		},
		{
			name:      "invalid UUID format",
			studentID: "invalid-uuid",
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:      "repository error",
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("DeleteStudent", mock.Anything, studentID).Return(errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:      "delete non-existent student", 
			studentID: studentID.String(),
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("DeleteStudent", mock.Anything, studentID).Return(errors.New("student not found"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Delete("/students/:id", handler.DeleteStudent)

			req := httptest.NewRequest("DELETE", "/students/"+tt.studentID, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			//////////////////////////////////////////////////
			// Response body validation
			//////////////////////////////////////////////////
			if !tt.wantErr && resp.StatusCode == fiber.StatusNoContent {
				// Success case - should have empty response body
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Empty(t, body)
			}

			if tt.wantErr {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				assert.NoError(t, err)
				assert.Contains(t, errorResp, "error")

				switch tt.name {
				case "invalid UUID format":
					assert.Contains(t, errorResp["error"], "Invalid UUID format")
				case "repository error", "delete non-existent student":
					assert.Contains(t, errorResp["error"], "Database error")
				}
			}
		})
	}
}