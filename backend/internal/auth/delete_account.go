package auth

import (
	"fmt"
	"io"
	"net/http"
	"specialstandard/internal/config"
	"specialstandard/internal/errs"
)

func SupabaseDeleteAccount(cfg *config.Supabase, userID string) error {
	supabaseURL := cfg.URL
	serviceRoleKey := cfg.ServiceRoleKey

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/auth/v1/admin/users/%s", supabaseURL, userID), nil)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", serviceRoleKey))
	req.Header.Set("apikey", serviceRoleKey)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("Failed to execute request: %v", err))
	}
	defer func() {
		_ = res.Body.Close()
	}()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return errs.BadRequest("failed to read response body")
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		return errs.BadRequest(fmt.Sprintf("failed to delete account, status: %d, response: %s", res.StatusCode, string(body)))
	}

	return nil
}
