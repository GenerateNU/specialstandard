package auth

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type Payload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserSignupResponse struct {
	ID uuid.UUID `json:"id"`
}

type SignupResponse struct {
	AccessToken string             `json:"access_token"`
	User        UserSignupResponse `json:"user"`
}

func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be atleast 8 characters long")
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#~$%^&*()+|_.,;<>?/{}\-]`).MatchString(password)

	if !hasUpper || !hasLower || !hasDigit || !hasSpecial {
		return errors.New("password must include uppercase, lowercase, digit and special characters")
	}

	return nil
}

func SupabaseSignup(cfg *config.Supabase, email, password string) (SignupResponse, error) {
	if err := validatePasswordStrength(password); err != nil {
		return SignupResponse{}, errs.BadRequest(fmt.Sprintf("Weak Password: %v", err))
	}

	supabaseURL := cfg.URL
	apiKey := cfg.ServiceRoleKey

	payload := Payload{
		Email:    email,
		Password: password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return SignupResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/signup", supabaseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		slog.Error("Error in Request Creation: ", "err", err)
		return SignupResponse{}, errs.BadRequest(fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("apikey", apiKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Error executing request: ", "err", err)
		return SignupResponse{}, errs.BadRequest(fmt.Sprintf("Failed to execute request: %v, %s", err, supabaseURL))
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Error reading response body: ", "body", body)
		return SignupResponse{}, errs.BadRequest("Failed to read response body", string(body))
	}

	if res.StatusCode != http.StatusOK {
		slog.Error("Error Response: ", "res.StatusCode", res.StatusCode, "body", string(body))
		return SignupResponse{}, errs.BadRequest(fmt.Sprintf("Failed to login %d, %s", res.StatusCode, body))
	}

	var response SignupResponse
	if err := json.Unmarshal(body, &response); err != nil {
		slog.Error("Error parsing response: ", "err", err)
		return SignupResponse{}, errs.BadRequest("Failed to parse request")
	}

	return response, nil
}
