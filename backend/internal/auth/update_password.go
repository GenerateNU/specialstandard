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

// Equivalent of Update-Password
func SupabaseUpdatePassword(cfg *config.Supabase, token, newPassword string) error {
	supabaseURL := cfg.URL
	apiKey := cfg.ServiceRoleKey

	payload := struct {
		Password string `json:"password"`
	}{
		Password: newPassword,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/auth/v1/user", supabaseURL), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("failed to create request: %v", err))
	}

	// Set Headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apiKey)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	res, err := Client.Do(req)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("failed to execute request: %v", err))
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errs.BadRequest("failed to read response body")
	}

	if res.StatusCode != http.StatusOK {
		return errs.BadRequest(fmt.Sprintf("failed to update password %d, %s", res.StatusCode, body))
	}

	return nil
}
