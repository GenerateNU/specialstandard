package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"

	"github.com/goccy/go-json"
)

func SupabaseForgotPassword(cfg *config.Supabase, email string, redirectURL string) error {
	supbaseURL := cfg.URL
	apiKey := cfg.ServiceRoleKey

	payload := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/recover", supbaseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Failed to create request: %v\n", err)
		return errs.BadRequest(fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	req.Header.Set("redirect_to", redirectURL)

	res, err := Client.Do(req)
	if err != nil {
		fmt.Printf("Failed to execute request: %v\n", err)
		return errs.BadRequest(fmt.Sprintf("failed to execute request: %v", err))
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v\n", err)
		return errs.BadRequest("failed to read response body")
	}

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Failed to initiate password reset: %d, %s\n", res.StatusCode, body)
		return errs.BadRequest(fmt.Sprintf("failed to initiate password reset %d, %s", res.StatusCode, body))
	}

	return nil
}
