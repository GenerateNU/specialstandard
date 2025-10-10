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
	"specialstandard/internal/models"
	"specialstandard/internal/service/handler/student"
	"specialstandard/internal/storage/mocks"
	"specialstandard/internal/utils"

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

func ptrInt(i int) *int {
	return &i
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
						Grade:       ptrInt(99),
						IEP:         ptrString("Test IEP"),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "", utils.NewPagination()).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty students list",
			url:  "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "", utils.NewPagination()).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "repository error",
			url:  "",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "", utils.NewPagination()).Return(nil, errors.New("database error"))
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
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "", utils.Pagination{Page: 2, Limit: 5}).Return([]models.Student{}, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with grade filter",
			url:  "?grade=5",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "John",
						LastName:  "Doe",
						Grade:     ptrInt(5),
					},
				}
				m.On("GetStudents", mock.Anything, ptrInt(5), uuid.Nil, "", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
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
				m.On("GetStudents", mock.Anything, (*int)(nil), therapistID, "", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
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
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "John", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "successful get students with all filters",
			url:  "?grade=5&therapist_id=123e4567-e89b-12d3-a456-426614174000&name=John&page=1&limit=5",
			mockSetup: func(m *mocks.MockStudentRepository) {
				therapistID := uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")
				students := []models.Student{
					{
						ID:          uuid.New(),
						FirstName:   "John",
						LastName:    "Doe",
						Grade:       ptrInt(5),
						TherapistID: therapistID,
					},
				}
				m.On("GetStudents", mock.Anything, ptrInt(5), therapistID, "John", utils.Pagination{Page: 1, Limit: 5}).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "empty results with filters",
			url:  "?grade=12&name=Nonexistent",
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudents", mock.Anything, ptrInt(12), uuid.Nil, "Nonexistent", mock.AnythingOfType("utils.Pagination")).Return([]models.Student{}, nil)
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
				m.On("GetStudents", mock.Anything, (*int)(nil), uuid.Nil, "JOHN", mock.AnythingOfType("utils.Pagination")).Return(students, nil)
			},
			expectedStatus: fiber.StatusOK,
			wantErr:        false,
		},
		{
			name: "filter by grade with pagination",
			url:  "?grade=5&page=2&limit=3",
			mockSetup: func(m *mocks.MockStudentRepository) {
				students := []models.Student{
					{
						ID:        uuid.New(),
						FirstName: "Student",
						LastName:  "Four",
						Grade:     ptrInt(5),
					},
				}
				m.On("GetStudents", mock.Anything, ptrInt(5), uuid.Nil, "", utils.Pagination{Page: 2, Limit: 3}).Return(students, nil)
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
						assert.Equal(t, 99, *students[0].Grade)
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
					Grade:       ptrInt(99),
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
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
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
					assert.Equal(t, 99, *student.Grade)
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
			requestBody: `{"grade": 5}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrInt(4),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrInt(5),
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
					Grade:       ptrInt(4),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrInt(4),
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
			requestBody: `{"first_name": "Updated", "last_name": "TestStudent", "grade": 5}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrInt(4),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Updated",
					LastName:    "TestStudent",
					Grade:       ptrInt(5),
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
					Grade:       ptrInt(4),
					IEP:         ptrString("Original IEP"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetStudent", mock.Anything, studentID).Return(existingStudent, nil)
				m.On("UpdateStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrInt(4),
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
			requestBody: `{"grade": 5}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				// No mock setup needed - UUID parsing fails before repository call
			},
			expectedStatus: fiber.StatusBadRequest,
			wantErr:        true,
		},
		{
			name:        "student not found",
			studentID:   studentID.String(),
			requestBody: `{"grade": 5}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("GetStudent", mock.Anything, studentID).Return(models.Student{}, errors.New("no rows in result set"))
			},
			expectedStatus: fiber.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:        "invalid JSON body",
			studentID:   studentID.String(),
			requestBody: `{"grade": 5 /* missing comma */}`,
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
					Grade:       ptrInt(4),
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
					Grade:       ptrInt(4),
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
			requestBody: `{"grade": 5}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				existingStudent := models.Student{
					ID:          studentID,
					FirstName:   "Test",
					LastName:    "Student",
					DOB:         ptrTime(time.Date(2011, 8, 12, 0, 0, 0, 0, time.UTC)),
					TherapistID: therapistID,
					Grade:       ptrInt(4),
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
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
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
				"grade": 5,
				"iep": "Active IEP with speech therapy goals"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "John",
					LastName:    "Doe",
					Grade:       ptrInt(5),
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
				"grade": 3,
				"iep": "Math accommodations and extended time"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "Emma",
					LastName:    "Johnson",
					Grade:       ptrInt(3),
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
				"grade": 12,
				"iep": "Graduation accommodations"
			}`,
			mockSetup: func(m *mocks.MockStudentRepository) {
				m.On("AddStudent", mock.Anything, mock.AnythingOfType("models.Student")).Return(models.Student{
					ID:          uuid.New(),
					FirstName:   "Test",
					LastName:    "Student",
					Grade:       ptrInt(12),
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
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
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
						assert.Equal(t, 5, *student.Grade)
					}
					if student.IEP != nil {
						assert.Contains(t, *student.IEP, "speech therapy")
					}
				case "successful create student with different data":
					assert.Equal(t, "Emma", student.FirstName)
					assert.Equal(t, "Johnson", student.LastName)
					if student.Grade != nil {
						assert.Equal(t, 3, *student.Grade)
					}
					if student.IEP != nil {
						assert.Contains(t, *student.IEP, "Math accommodations")
					}
				case "valid date edge cases":
					assert.Equal(t, "Test", student.FirstName)
					if student.Grade != nil {
						assert.Equal(t, 12, *student.Grade)
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

				// Check for HTTPError structure with code and message
				assert.Contains(t, errorResp, "code")
				assert.Contains(t, errorResp, "message")

				// Validate specific error messages based on test case
				switch tt.name {
				case "invalid date format":
					// For validation errors, message is a map
					if msgMap, ok := errorResp["message"].(map[string]interface{}); ok {
						assert.Contains(t, msgMap, "dob")
					}
				case "invalid therapist UUID format":
					// For validation errors, message is a map
					if msgMap, ok := errorResp["message"].(map[string]interface{}); ok {
						assert.Contains(t, msgMap, "therapistid")
					}
				case "missing required fields":
					// For validation errors, message is a map
					if msgMap, ok := errorResp["message"].(map[string]interface{}); ok {
						// Should have errors for missing fields
						assert.True(t, len(msgMap) > 0)
					}
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
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
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
