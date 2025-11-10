package auth

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/mocks"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_SignUp(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockTherapistRepository)
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name: "Invalid Request Body",
			payload: fmt.Sprintf(`{
				"id": "%s",
				"first_name": 123,
				"last_name": true,
				"email": "doctor.guess.who.suess@gmail.com"
			}`, uuid.New()),
			mockSetup:          func(m *mocks.MockTherapistRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
			wantErr:            true,
		},
		{
			name: "Successful Signup Request",
			payload: `{
				"email": "meow.thegato@gmail.com",
				"password": "Meow123;TunaToMe",
				"first_name": "El",
				"last_name": "Catto"
			}`,
			mockSetup: func(m *mocks.MockTherapistRepository) {
				therapist := &models.Therapist{
					ID:        uuid.New(),
					FirstName: "El",
					LastName:  "Catto",
					Email:     "meow.thegato@gmail.com",
					Active:    true,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				m.On("CreateTherapist", mock.Anything, mock.AnythingOfType("*models.CreateTherapistInput")).Return(therapist, nil)
			},
			expectedStatusCode: fiber.StatusCreated,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"access_token": "dummy-token",
					"user": {"id": "f20e5948-01ba-4113-b453-db05d8bde3bc"}
				}`))
			}))
			defer ts.Close()
			mockConfig := config.Supabase{
				URL:            ts.URL,
				ServiceRoleKey: "SRK",
			}

			handler := NewHandler(mockConfig, mockRepo)
			app.Post("/signup", handler.SignUp)

			req := httptest.NewRequest("POST", "/signup", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_Login(t *testing.T) {
	tests := []struct {
		name               string
		payload            string
		mockSetup          func(*mocks.MockTherapistRepository)
		expectedStatusCode int
		wantErr            bool
	}{
		{
			name: "Invalid Request Body",
			payload: `{
				"email": 123,
				"password": true
			}`,
			mockSetup:          func(m *mocks.MockTherapistRepository) {},
			expectedStatusCode: fiber.StatusBadRequest,
			wantErr:            true,
		},
		{
			name: "Successful Login Request",
			payload: `{
				"email": "meow.thegato@gmail.com",
				"password": "Meow123;TunaToMe"
			}`,
			mockSetup:          func(m *mocks.MockTherapistRepository) {},
			expectedStatusCode: fiber.StatusOK,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New(fiber.Config{
				ErrorHandler: errs.ErrorHandler,
			})
			mockRepo := new(mocks.MockTherapistRepository)
			tt.mockSetup(mockRepo)

			// Test Supabase server for login
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{
					"access_token": "dummy-token",
					"user": {"id": "f20e5948-01ba-4113-b453-db05d8bde3bc"}
				}`))
			}))
			defer ts.Close()

			mockConfig := config.Supabase{
				URL:            ts.URL,
				ServiceRoleKey: "SRK",
			}

			handler := NewHandler(mockConfig, mockRepo)
			app.Post("/login", handler.Login)

			req := httptest.NewRequest("POST", "/login", strings.NewReader(tt.payload))
			req.Header.Set("Content-Type", "application/json")
			res, _ := app.Test(req, -1)

			assert.Equal(t, tt.expectedStatusCode, res.StatusCode)
			mockRepo.AssertExpectations(t)
		})
	}
}
