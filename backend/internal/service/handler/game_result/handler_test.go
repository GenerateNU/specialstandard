package game_result

import (
	"errors"
	"fmt"
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/mocks"
	"strings"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetGameResults(t *testing.T) {
	studentID := uuid.MustParse("f4053a11-fc76-4551-b05e-b550dfea516d")
	sessionID := uuid.MustParse("b6051b20-5426-428e-858e-adbe853244e3")

	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockGameResultRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "Empty DB Test",
			url:  "",
			mockSetup: func(m *mocks.MockGameResultRepository) {
				m.On("GetGameResults", mock.Anything, mock.AnythingOfType("*models.GetGameResultQuery"), mock.AnythingOfType("utils.Pagination")).Return([]models.GameResult{}, nil)
			},
			expectedStatus: 200,
			wantErr:        false,
		},
		{
			name:           "Parsing Error - Bad Student ID",
			url:            "?student_id=123",
			mockSetup:      func(m *mocks.MockGameResultRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name:           "Invalid Query-Param, Validation Issue, Page: -1",
			url:            "?page=-1",
			mockSetup:      func(m *mocks.MockGameResultRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name:           "Invalid Query-Param, Parsing Issue, Page: abc",
			url:            "?page=abc",
			mockSetup:      func(m *mocks.MockGameResultRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name: "Valid Case",
			url:  "",
			mockSetup: func(m *mocks.MockGameResultRepository) {
				gameResult := models.GameResult{
					ID:             uuid.New(),
					SessionID:      uuid.New(),
					StudentID:      uuid.New(),
					ContentID:      uuid.New(),
					TimeTaken:      40,
					Completed:      true,
					IncorrectTries: 3,
					CreatedAt:      ptr.Time(time.Now()),
					UpdatedAt:      ptr.Time(time.Now()),
				}

				m.On("GetGameResults", mock.Anything, mock.AnythingOfType("*models.GetGameResultQuery"), mock.AnythingOfType("utils.Pagination")).Return([]models.GameResult{
					gameResult, gameResult, gameResult,
				}, nil)
			},
			expectedStatus: 200,
			wantErr:        false,
		},
		{
			name: "Valid Case",
			url:  "?student_id=f4053a11-fc76-4551-b05e-b550dfea516d",
			mockSetup: func(m *mocks.MockGameResultRepository) {
				gameResult := models.GameResult{
					ID:             uuid.New(),
					SessionID:      uuid.New(),
					StudentID:      studentID,
					ContentID:      uuid.New(),
					TimeTaken:      40,
					Completed:      true,
					IncorrectTries: 3,
					CreatedAt:      ptr.Time(time.Now()),
					UpdatedAt:      ptr.Time(time.Now()),
				}

				m.On("GetGameResults", mock.Anything, mock.AnythingOfType("*models.GetGameResultQuery"), mock.AnythingOfType("utils.Pagination")).Return([]models.GameResult{
					gameResult,
				}, nil)
			},
			expectedStatus: 200,
			wantErr:        false,
		},
		{
			name: "Valid Case",
			url:  "?session_id=b6051b20-5426-428e-858e-adbe853244e3",
			mockSetup: func(m *mocks.MockGameResultRepository) {
				gameResult := models.GameResult{
					ID:             uuid.New(),
					SessionID:      sessionID,
					StudentID:      uuid.New(),
					ContentID:      uuid.New(),
					TimeTaken:      40,
					Completed:      true,
					IncorrectTries: 3,
					CreatedAt:      ptr.Time(time.Now()),
					UpdatedAt:      ptr.Time(time.Now()),
				}

				m.On("GetGameResults", mock.Anything, mock.AnythingOfType("*models.GetGameResultQuery"), mock.AnythingOfType("utils.Pagination")).Return([]models.GameResult{
					gameResult,
				}, nil)
			},
			expectedStatus: 200,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockGameResultRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			app.Get("/game-results", handler.GetGameResults)

			req := httptest.NewRequest("GET", "/game-results"+tt.url, nil)
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_PostGameResult(t *testing.T) {
	sessionID := uuid.New()
	studentID := uuid.New()
	contentID := uuid.New()

	tests := []struct {
		name           string
		payload        string
		mockSetup      func(*mocks.MockGameResultRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "Parsing Error",
			payload: fmt.Sprintf(`{
				"session_id": "%s",
				"student_id": "%s",
				"content_id": "%s",
				"time_taken": 93,
				"completed": true,
				"incorrect_tries": 1,
            }`, sessionID, studentID, contentID),
			mockSetup:      func(m *mocks.MockGameResultRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name: "Invalid Time-Taken",
			payload: fmt.Sprintf(`{
				"session_id": "%s",
				"student_id": "%s",
				"content_id": "%s",
				"time_taken": -93,
				"completed": true,
				"incorrect_tries": 1
            }`, sessionID, studentID, contentID),
			mockSetup:      func(m *mocks.MockGameResultRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name: "Foreign Key Reference Error",
			payload: fmt.Sprintf(`{
				"session_id": "%s",
				"student_id": "%s",
				"content_id": "%s",
				"time_taken": 93,
				"completed": true,
				"incorrect_tries": 1
            }`, sessionID, studentID, contentID),
			mockSetup: func(m *mocks.MockGameResultRepository) {
				m.On("PostGameResult", mock.Anything, mock.AnythingOfType("PostGameResult")).Return(nil, errors.New("foreign key"))
			},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name: "Valid without Optional Parameter",
			payload: fmt.Sprintf(`{
				"session_id": "%s",
				"student_id": "%s",
				"content_id": "%s",
				"time_taken": 93,
				"completed": false
            }`, sessionID, studentID, contentID),
			mockSetup: func(m *mocks.MockGameResultRepository) {
				postGameResult := models.PostGameResult{
					SessionID: sessionID,
					StudentID: studentID,
					ContentID: contentID,
					TimeTaken: 93,
					Completed: ptr.Bool(false),
				}

				gameResult := &models.GameResult{
					ID:             uuid.New(),
					SessionID:      sessionID,
					StudentID:      studentID,
					ContentID:      contentID,
					TimeTaken:      93,
					Completed:      false,
					IncorrectTries: 0,
					CreatedAt:      ptr.Time(time.Now()),
					UpdatedAt:      ptr.Time(time.Now()),
				}
				m.On("PostGameResult", mock.Anything, postGameResult).Return(gameResult, nil)
			},
			expectedStatus: 201,
			wantErr:        false,
		},
		{
			name: "Valid with Optional Parameters",
			payload: fmt.Sprintf(`{
				"session_id": "%s",
				"student_id": "%s",
				"content_id": "%s",
				"time_taken": 93,
				"completed": false,
				"incorrect_tries": 9
            }`, sessionID, studentID, contentID),
			mockSetup: func(m *mocks.MockGameResultRepository) {
				postGameResult := models.PostGameResult{
					SessionID:      sessionID,
					StudentID:      studentID,
					ContentID:      contentID,
					TimeTaken:      93,
					Completed:      ptr.Bool(false),
					IncorrectTries: ptr.Int(9),
				}

				gameResult := &models.GameResult{
					ID:             uuid.New(),
					SessionID:      sessionID,
					StudentID:      studentID,
					ContentID:      contentID,
					TimeTaken:      93,
					Completed:      false,
					IncorrectTries: 9,
					CreatedAt:      ptr.Time(time.Now()),
					UpdatedAt:      ptr.Time(time.Now()),
				}
				m.On("PostGameResult", mock.Anything, postGameResult).Return(gameResult, nil)
			},
			expectedStatus: 201,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockGameResultRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			app.Post("/game-results", handler.PostGameResult)

			req := httptest.NewRequest("POST", "/game-results", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")

			res, _ := app.Test(req, -1)
			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
