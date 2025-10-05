package student_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/student"
	"specialstandard/internal/storage/mocks"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrTime(t time.Time) *time.Time {
	return &t
}

func ptrString(s string) *string {
	return &s
}

func TestHandler_GetStudents(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "successful get students with default pagination",
			url:  "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:          uuid.New(),
						FirstName:   "Test",
						LastName:    "Student",
						DOB:         ptrTime(time.Now().AddDate(-10, 0, 0)),
						TherapistID: uuid.New(),
						Grade:       ptrString("Test Grade"),
						IEP:         ptrString("Test IEP"),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "", utils.NewPagination()).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty students list",
			url:  "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "", utils.NewPagination()).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "", utils.NewPagination()).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest, // QueryParser Fails
			wantErr:        true,
		},
		{
			name: "Pagination Parameters",
			url:  "?page=2&limit=5",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "", utils.Pagination{Page: 2, Limit: 5}).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with grade filter",
			url:  "?grade=5th",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "John",
						LastName:  "Doe",
						Grade:     ptrString("5th"),
					},
				}
				m.On("GetStudents", mock.Anything, "5th", uuid.Nil, "", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with therapist filter",
			url:  "?therapist_id=123e4567-e89b-12d3-a456-426614174000",
			mockSetup: func(m *mocks.MockStudentRepository) {
				therapistID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				students := []models.Student{
					{
						ID:          uuid.New(),
						FirstName:   "Jane",
						LastName:    "Smith",
						TherapistID: therapistID,
					},
				}
				m.On("GetStudents", mock.Anything, "", therapistID, "", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with name filter",
			url:  "?name=John",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "John",
						LastName:  "Doe",
					},
				}
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "John", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with all filters",
			url:  "?grade=5th&therapist_id=123e4567-e89b-12d3-a456-426614174000&name=John&page=1&limit=5",
			mockSetup: func(m *mocks.MockStudentRepository) {
				therapistID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				students := []models.Student{
					{
						ID:          uuid.New(),
						FirstName:   "John",
						LastName:    "Doe",
						Grade:       ptrString("5th"),
						TherapistID: therapistID,
					},
				}
				m.On("GetStudents", mock.Anything, "5th", therapistID, "John", utils.Pagination{Page: 1, Limit: 5}).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty results with filters",
			url:  "?grade=12th&name=Nonexistent",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, "12th", uuid.Nil, "Nonexistent", mock.AnythingOfType("utils.Pagination")).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "case insensitive name search",
			url:  "?name=JOHN",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "John",
						LastName:  "Doe",
					},
				}
				m.On("GetStudents", mock.Anything, "", uuid.Nil, "JOHN", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "filter by grade with pagination",
			url:  "?grade=5th&page=2&limit=3",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "Student",
						LastName:  "Four",
						Grade:     ptrString("5th"),
					},
				}
				m.On("GetStudents", mock.Anything, "5th", uuid.Nil, "", utils.Pagination{Page: 2, Limit: 3}).Return(students, nil)
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
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Get("/students", handler.GetStudents)

			req := httptest.NewRequest("GET", "/students"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			// Response body validation
			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var students []models.Student
				err = json.Unmarshal(body, &students)
				assert.NoError(t, err)

				switch tt.name {
				case "empty students list":
					assert.Len(t, students, 0)
				case "successful get students":
					assert.Len(t, students, 1)
					assert.Equal(t, "Test", students[0].FirstName)
					assert.Equal(t, "Student", students[0].LastName)
					// Update assertions to handle nullable pointers
					if students[0].Grade != nil {
						assert.Equal(t, "Test Grade", *students[0].Grade)
					}
					if students[0].IEP != nil {
						assert.Equal(t, "Test IEP", *students[0].IEP)
					}
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
					DOB:         ptrTime(time.Now().AddDate(-10, 0, 0)),
					TherapistID: uuid.New(),
					Grade:       ptrString("Test Grade"),
					IEP:         ptrString("Test IEP"),
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

			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				// Success case - validate student data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var student models.Student
				err = json.Unmarshal(body, &student)
				assert.NoError(t, err)

				// Validate the student data with pointer handling
				assert.Equal(t, "Test", student.FirstName)
				assert.Equal(t, "Student", student.LastName)
				if student.Grade != nil {
					assert.Equal(t, "Test Grade", *student.Grade)
				}
				if student.IEP != nil {
					assert.Equal(t, "Test IEP", *student.IEP)
				}
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
		{
			name:        "update grade only",
			studentID:   studentID.String(),
			requestBody: `{"grade": "5th"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrString("5th"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "update IEP only",
			studentID:   studentID.String(),
			requestBody: `{"iep": "Updated IEP with math accommodations"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrString("4th"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					IEP:         ptrString("Updated IEP with math accommodations"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "update name and grade",
			studentID:   studentID.String(),
			requestBody: `{"first_name": "Updated", "last_name": "TestStudent", "grade": "5th"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Updated",
					LastName:    "TestStudent",
					Grade:       ptrString("5th"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "update DOB with valid date",
			studentID:   studentID.String(),
			requestBody: `{"dob": "2010-05-15"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrString("4th"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:        "invalid UUID format",
			studentID:   "invalid-uuid",
			requestBody: `{"grade": "5th"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "student not found",
			studentID:   studentID.String(),
			requestBody: `{"grade": "5th"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudent", mock.Anything, studentID).Return(models.Student{}, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:        "invalid JSON body",
			studentID:   studentID.String(),
			requestBody: `{"grade": "5th" /* missing comma */}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid date format",
			studentID:   studentID.String(),
			requestBody: `{"dob": "invalid-date"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "invalid therapist UUID",
			studentID:   studentID.String(),
			requestBody: `{"therapist_id": "bad-uuid"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "UpdateStudent repository error",
			studentID:   studentID.String(),
			requestBody: `{"grade": "5th"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrString("4th"),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{}, errors.New("database error"))
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

			// Make request
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
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "John",
					LastName:    "Doe",
					Grade:       ptrString("5th"),
					IEP:         ptrString("Active IEP with speech therapy goals"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2010, 5, 15, 0, 0, 0, 0, time.UTC)),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
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
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "Emma",
					LastName:    "Johnson",
					Grade:       ptrString("3rd"),
					IEP:         ptrString("Math accommodations and extended time"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2012, 3, 22, 0, 0, 0, 0, time.UTC)),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
			},
			expectedStatus: fiber.StatusCreated,
			wantErr:        false,
		},
		{
			name:        "invalid JSON body",
			requestBody: `{"first_name": "John", "last_name": "Doe" /* missing comma */}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - JSON parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "missing required fields",
			requestBody: `{"first_name": "John"}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - validation fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "invalid date format",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"dob": "invalid-date",
				"therapist_id": "` + therapistID.String() + `"
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
				"therapist_id": "invalid-uuid"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "repository save error",
			requestBody: `{
				"first_name": "John",
				"last_name": "Doe",
				"therapist_id": "` + therapistID.String() + `"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{}, errors.New("database connection failed"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
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
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrString("12th"),
					TherapistID: therapistID,
					DOB:         ptrTime(time.Date(2000, 2, 29, 0, 0, 0, 0, time.UTC)),
					IEP:         ptrString("Graduation accommodations"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil)
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

			if !tt.wantErr && resp.StatusCode == fiber.StatusCreated {
				// Success case - validate created student data
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var student models.Student
				err = json.Unmarshal(body, &student)
				assert.NoError(t, err)

				// Validate response data with pointer handling
				switch tt.name {
				case "successful create student":
					assert.Equal(t, "John", student.FirstName)
					assert.Equal(t, "Doe", student.LastName)
					if student.Grade != nil {
						assert.Equal(t, "5th", *student.Grade)
					}
					if student.IEP != nil {
						assert.Contains(t, *student.IEP, "speech therapy")
					}
				case "successful create student with different data":
					assert.Equal(t, "Emma", student.FirstName)
					assert.Equal(t, "Johnson", student.LastName)
					if student.Grade != nil {
						assert.Equal(t, "3rd", *student.Grade)
					}
					if student.IEP != nil {
						assert.Contains(t, *student.IEP, "Math accommodations")
					}
				case "valid date edge cases":
					assert.Equal(t, "Test", student.FirstName)
					if student.Grade != nil {
						assert.Equal(t, "12th", *student.Grade)
					}
				}

				// Validate that UUID was generated
				assert.NotEqual(t, uuid.Nil, student.ID)
				assert.Equal(t, therapistID, student.TherapistID)

				// Validate date was parsed correctly if provided
				if student.DOB != nil {
					assert.False(t, student.DOB.IsZero())
				}
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
					assert.Contains(t, errorResp["error"], "required")
				}
			}
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
				case "repository error":
					assert.Contains(t, errorResp["error"], "Database error")
				}
			}
		})
	}
}

func TestHandler_GetStudentSessions(t *testing.T) {
	studentID := uuid.New()
	sessionID := uuid.New()
	therapistID := uuid.New()
	
	// Helper to create mock session data
	createMockSession := func(startTime time.Time, present bool) models.StudentSessionsOutput {
		return models.StudentSessionsOutput{
			StudentID: studentID,
			Present:   present,
			Notes:     ptrString("Test session notes"),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Session: models.Session{
				ID:            sessionID,
				StartDateTime: startTime,
				EndDateTime:   startTime.Add(time.Hour),
				TherapistID:   therapistID,
				Notes:         ptrString("Session notes"),
				CreatedAt:     ptrTime(time.Now()),
				UpdatedAt:     ptrTime(time.Now()),
			},
		}
	}

	tests := []struct {
		name           string
		studentID      string
		url            string
		mockSetup      func(*mocks.MockStudentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name:      "successful get student sessions with default pagination",
			studentID: studentID.String(),
			url:       "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Now(), true),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), (*models.GetStudentSessionsRepositoryRequest)(nil)).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "empty sessions list",
			studentID: studentID.String(),
			url:       "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), (*models.GetStudentSessionsRepositoryRequest)(nil)).Return([]models.StudentSessionsOutput{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "repository error",
			studentID: studentID.String(),
			url:       "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), (*models.GetStudentSessionsRepositoryRequest)(nil)).Return(nil, errors.New("database error"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:      "invalid UUID format",
			studentID: "invalid-uuid",
			url:       "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		// ------- Pagination Cases -------
		{
			name:           "Violating Pagination Arguments Constraints",
			studentID:      studentID.String(),
			url:            "?page=0&limit=-1",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Bad Pagination Arguments",
			studentID:      studentID.String(),
			url:            "?page=abc&limit=-1",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:      "Pagination Parameters",
			studentID: studentID.String(),
			url:       "?page=2&limit=5",
			mockSetup: func(m *mocks.MockStudentRepository) {
				pagination := utils.Pagination{Page: 2, Limit: 5}
				m.On("GetStudentSessions", mock.Anything, studentID, pagination, (*models.GetStudentSessionsRepositoryRequest)(nil)).Return([]models.StudentSessionsOutput{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		// ------- Attendance Filter Cases -------
		{
			name:      "Filter by present=true",
			studentID: studentID.String(),
			url:       "?present=true",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Now(), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Present: func() *bool { b := true; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by present=false",
			studentID: studentID.String(),
			url:       "?present=false",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Now(), false),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Present: func() *bool { b := false; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		// ------- Date Range Filter Cases -------
		{
			name:      "Filter by date range",
			studentID: studentID.String(),
			url:       "?startDate=2025-09-01&endDate=2025-09-30",
			mockSetup: func(m *mocks.MockStudentRepository) {
				startDate, _ := time.Parse("2006-01-02", "2025-09-01")
				endDate, _ := time.Parse("2006-01-02", "2025-09-30")
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					StartDate: &startDate,
					EndDate:   &endDate,
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by start date only",
			studentID: studentID.String(),
			url:       "?startDate=2025-09-01",
			mockSetup: func(m *mocks.MockStudentRepository) {
				startDate, _ := time.Parse("2006-01-02", "2025-09-01")
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					StartDate: &startDate,
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by end date only",
			studentID: studentID.String(),
			url:       "?endDate=2025-09-30",
			mockSetup: func(m *mocks.MockStudentRepository) {
				endDate, _ := time.Parse("2006-01-02", "2025-09-30")
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					EndDate: &endDate,
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		// ------- Month/Year Filter Cases -------
		{
			name:      "Filter by month and year",
			studentID: studentID.String(),
			url:       "?month=9&year=2025",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Month: func() *int { i := 9; return &i }(),
					Year:  func() *int { i := 2025; return &i }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by year only",
			studentID: studentID.String(),
			url:       "?year=2025",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Year: func() *int { i := 2025; return &i }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by month only (should work with current year)",
			studentID: studentID.String(),
			url:       "?month=9",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Month: func() *int { i := 9; return &i }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		// ------- Combined Filter Cases -------
		{
			name:      "Filter by attendance and date range",
			studentID: studentID.String(),
			url:       "?present=true&startDate=2025-09-01&endDate=2025-09-30",
			mockSetup: func(m *mocks.MockStudentRepository) {
				startDate, _ := time.Parse("2006-01-02", "2025-09-01")
				endDate, _ := time.Parse("2006-01-02", "2025-09-30")
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					StartDate: &startDate,
					EndDate:   &endDate,
					Present:   func() *bool { b := true; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "Filter by attendance and month/year",
			studentID: studentID.String(),
			url:       "?present=false&month=9&year=2025",
			mockSetup: func(m *mocks.MockStudentRepository) {
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), false),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Month:   func() *int { i := 9; return &i }(),
					Year:    func() *int { i := 2025; return &i }(),
					Present: func() *bool { b := false; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name:      "All filters combined with pagination",
			studentID: studentID.String(),
			url:       "?present=true&month=9&year=2025&page=2&limit=3",
			mockSetup: func(m *mocks.MockStudentRepository) {
				pagination := utils.Pagination{Page: 2, Limit: 3}
				sessions := []models.StudentSessionsOutput{
					createMockSession(time.Date(2025, 9, 15, 10, 0, 0, 0, time.UTC), true),
				}
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Month:   func() *int { i := 9; return &i }(),
					Year:    func() *int { i := 2025; return &i }(),
					Present: func() *bool { b := true; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, pagination, expectedFilter).Return(sessions, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		// ------- Invalid Query Parameter Cases -------
		{
			name:           "Invalid month value (too low)",
			studentID:      studentID.String(),
			url:            "?month=0",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid month value (too high)",
			studentID:      studentID.String(),
			url:            "?month=13",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid year value (too low)",
			studentID:      studentID.String(),
			url:            "?year=1775",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid year value (too high)",
			studentID:      studentID.String(),
			url:            "?year=2201",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid present boolean value",
			studentID:      studentID.String(),
			url:            "?present=maybe",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid startDate format",
			studentID:      studentID.String(),
			url:            "?startDate=invalid-date",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:           "Invalid endDate format",
			studentID:      studentID.String(),
			url:            "?endDate=2025-13-45",
			mockSetup:      func(m *mocks.MockStudentRepository) {},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		// ------- Empty Filter Results Cases -------
		{
			name:      "No sessions match filters",
			studentID: studentID.String(),
			url:       "?present=true&year=2024",
			mockSetup: func(m *mocks.MockStudentRepository) {
				expectedFilter := &models.GetStudentSessionsRepositoryRequest{
					Year:    func() *int { i := 2024; return &i }(),
					Present: func() *bool { b := true; return &b }(),
				}
				m.On("GetStudentSessions", mock.Anything, studentID, utils.NewPagination(), expectedFilter).Return([]models.StudentSessionsOutput{}, nil)
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
			mockRepo := new(mocks.MockStudentRepository)
			tt.mockSetup(mockRepo)

			handler := student.NewHandler(mockRepo)
			app.Get("/students/:id/sessions", handler.GetStudentSessions)

			req := httptest.NewRequest("GET", "/students/"+tt.studentID+"/sessions"+tt.url, nil)
			resp, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			mockRepo.AssertExpectations(t)

			// Response body validation for successful cases
			if !tt.wantErr && resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var sessions []models.StudentSessionsOutput
				err = json.Unmarshal(body, &sessions)
				assert.NoError(t, err)

				// Validate response structure
				if len(sessions) > 0 {
					session := sessions[0]
					assert.Equal(t, studentID, session.StudentID)
					assert.NotNil(t, session.Session.ID)
					assert.False(t, session.Session.StartDateTime.IsZero())
					assert.False(t, session.Session.EndDateTime.IsZero())
					
					// Validate filter-specific expectations
					switch {
					case strings.Contains(tt.url, "present=true"):
						assert.True(t, session.Present)
					case strings.Contains(tt.url, "present=false"):
						assert.False(t, session.Present)
					}
				}
			}

			// Error response validation
			if tt.wantErr {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)

				var errorResp map[string]interface{}
				err = json.Unmarshal(body, &errorResp)
				assert.NoError(t, err)
				
				// The API returns {"code": X, "message": "..."} not {"error": "..."}
				assert.Contains(t, errorResp, "message")

				// Validate specific error messages based on status code and test name
				switch {
				case strings.Contains(tt.name, "Invalid UUID"):
					assert.Contains(t, errorResp["message"], "Invalid UUID format")
				case strings.Contains(tt.name, "repository error"):
					assert.Contains(t, errorResp["message"], "Failed to retrieve student sessions")
				case strings.Contains(tt.name, "month") && strings.Contains(tt.name, "Invalid"):
					assert.Contains(t, errorResp["message"], "Month must be between 1 and 12")
				case strings.Contains(tt.name, "year") && strings.Contains(tt.name, "Invalid"):
					assert.Contains(t, errorResp["message"], "Year must be between 1776 and 2200")
				case strings.Contains(tt.name, "present") && strings.Contains(tt.name, "Invalid"):
					assert.Contains(t, errorResp["message"], "Present must be 'true' or 'false'")
				case strings.Contains(tt.name, "Date") && strings.Contains(tt.name, "Invalid"):
					assert.Contains(t, errorResp["message"], "Invalid")
				case strings.Contains(tt.name, "Pagination"):
					// Pagination errors can have different formats
					assert.True(t, strings.Contains(fmt.Sprintf("%v", errorResp["message"]), "Pagination") || 
						strings.Contains(fmt.Sprintf("%v", errorResp["message"]), "limit") ||
						strings.Contains(fmt.Sprintf("%v", errorResp["message"]), "page"))
				}
			}
		})
	}
}
