package auth

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type userResponse struct {
	ID uuid.UUID `json:"id"`
}

type SignInResponse struct {
	AccessToken  string       `json:"access_token"`
	TokenType    string       `json:"token_type"`
	ExpiresIn    int          `json:"expires_in"`
	RefreshToken string       `json:"refresh_token"`
	User         userResponse `json:"user"`
	Error        interface{}  `json:"error"`
}

func SupabaseLogin(cfg *config.Supabase, email string, password string) (SignInResponse, error) {
	supabaseURL := cfg.URL
	serviceroleKey := cfg.ServiceRoleKey

	payload := Payload{
		Email:    email,
		Password: password,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return SignInResponse{}, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/token?grant_type=password", supabaseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return SignInResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceroleKey))
	req.Header.Set("apikey", serviceroleKey)

	res, err := Client.Do(req)
	if err != nil {
		slog.Error("Failed to execute Request", "err", err)
		return SignInResponse{}, errs.BadRequest("Failed to execute Request")
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slog.Error("Failed to read response body", "err", err)
		return SignInResponse{}, errs.BadRequest("Failed to read response body")
	}

	if res.StatusCode != http.StatusOK {
    // Try to parse the error response
    var errorResp struct {
        Message string `json:"msg"`
        Error   string `json:"error"`
    }
    
    if err := json.Unmarshal(body, &errorResp); err == nil {
        // Use the parsed message if available
        if errorResp.Message != "" {
            return SignInResponse{}, errs.BadRequest(errorResp.Message)
        }
        if errorResp.Error != "" {
            return SignInResponse{}, errs.BadRequest(errorResp.Error)
        }
    }
    
    // Fallback to generic message if parsing fails
    return SignInResponse{}, errs.BadRequest("Invalid credentials")
}

	var signInResponse SignInResponse
	err = json.Unmarshal(body, &signInResponse)
	if err != nil {
		slog.Error("Failed to parse response body", "body", err)
		return SignInResponse{}, errs.BadRequest("Failed to parse response body")
	}

	if signInResponse.Error != nil {
		return SignInResponse{}, errs.BadRequest(fmt.Sprintf("Sign In Response Error %v", signInResponse.Error))
	}

	return signInResponse, nil
}
