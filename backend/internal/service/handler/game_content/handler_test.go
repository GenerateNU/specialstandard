package game_content

import (
	"net/http/httptest"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/mocks"
	"testing"
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetGameContents(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		mockSetup      func(*mocks.MockGameContentRepository)
		expectedStatus int
		wantErr        bool
	}{
		{
			name: "404 Not Found Error",
			url:  "?category=following_directions&level=4&count=5",
			mockSetup: func(m *mocks.MockGameContentRepository) {
				m.On("GetGameContent", mock.Anything, mock.AnythingOfType("models.GetGameContentRequest")).Return(nil, pgx.ErrNoRows)
			},
			expectedStatus: 404,
			wantErr:        true,
		},
		{
			name: "Overboard word-count returns all options",
			url:  "?category=sequencing&level=5&count=100",
			mockSetup: func(m *mocks.MockGameContentRepository) {
				gameContent := &models.GameContent{
					ID:              uuid.New(),
					Category:        "sequencing",
					DifficultyLevel: 5,
					Options:         []string{"GeneRAT", "Meow", "Liepard", "Cat", "Dog", "Oink"},
					Answer:          "Woof",
					CreatedAt:       ptr.Time(time.Now()),
					UpdatedAt:       ptr.Time(time.Now()),
				}
				m.On("GetGameContent", mock.Anything, mock.AnythingOfType("models.GetGameContentRequest")).Return(gameContent, nil)
			},
			expectedStatus: 200,
			wantErr:        false,
		},
		{
			name:           "Invalid category enum",
			url:            "?category=Sequencing&level=5&count=100",
			mockSetup:      func(m *mocks.MockGameContentRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name:           "Missing DifficultyLevel-Field in Query",
			url:            "?category=sequencing&count=100",
			mockSetup:      func(m *mocks.MockGameContentRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name:           "Invalid Count Query-Parameter",
			url:            "?category=sequencing&level=5&count=1",
			mockSetup:      func(m *mocks.MockGameContentRepository) {},
			expectedStatus: 400,
			wantErr:        true,
		},
		{
			name: "Valid Case",
			url:  "?category=sequencing&level=5&count=3",
			mockSetup: func(m *mocks.MockGameContentRepository) {
				gameContent := &models.GameContent{
					ID:              uuid.New(),
					Category:        "sequencing",
					DifficultyLevel: 5,
					Options:         []string{"GeneRAT", "Meow"},
					Answer:          "Woof",
					CreatedAt:       ptr.Time(time.Now()),
					UpdatedAt:       ptr.Time(time.Now()),
				}
				m.On("GetGameContent", mock.Anything, mock.AnythingOfType("models.GetGameContentRequest")).Return(gameContent, nil)
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
			mockRepo := new(mocks.MockGameContentRepository)
			tt.mockSetup(mockRepo)

			handler := NewHandler(mockRepo)
			app.Get("/game-contents", handler.GetGameContents)

			req := httptest.NewRequest("GET", "/game-contents"+tt.url, nil)
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatus, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
