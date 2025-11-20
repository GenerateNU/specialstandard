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

func SupabaseUpdatePassword(cfg *config.Supabase, email, recoveryToken, newPassword string) error {
	supabaseURL := cfg.URL
	apiKey := cfg.ServiceRoleKey

	fmt.Printf("SupabaseUpdatePassword called with email: %s\n", email)
	fmt.Printf("SupabaseUpdatePassword called with token: %s\n", recoveryToken)
	fmt.Printf("SupabaseUpdatePassword called with new password: %s\n", newPassword)

	// Step 1: Verify recovery token with email
	verifyPayload := map[string]string{
		"type":  "recovery",
		"email": email,
		"token": recoveryToken,
	}
	fmt.Printf("SupabaseUpdatePassword called with email: %s, token: %s\n", email, recoveryToken)

	verifyBytes, _ := json.Marshal(verifyPayload)
	verifyReq, _ := http.NewRequest("POST", fmt.Sprintf("%s/auth/v1/verify", supabaseURL), bytes.NewBuffer(verifyBytes))
	verifyReq.Header.Set("Content-Type", "application/json")
	verifyReq.Header.Set("apikey", apiKey)

	res, err := Client.Do(verifyReq)
	if err != nil {
		fmt.Printf("Failed to verify recovery token request: %v\n", err)
		return errs.BadRequest("Failed to verify recovery token")
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		fmt.Printf("Invalid or expired recovery token response: %d, %s\n", res.StatusCode, string(body))
		return errs.BadRequest(fmt.Sprintf("Invalid or expired recovery token: %s", string(body)))
	}

	var verifyResp struct {
		Session struct {
			AccessToken string `json:"access_token"`
		} `json:"session"`
	}

	fmt.Printf("Verify response body: %s\n", string(body))

	if err := json.Unmarshal(body, &verifyResp); err != nil {
		fmt.Printf("Failed to parse token response: %v\n", err)
		return errs.BadRequest("Failed to parse token response")
	}

	fmt.Printf("Access token: %s\n", verifyResp.Session.AccessToken)

	accessToken := verifyResp.Session.AccessToken
	if accessToken == "" {
		fmt.Printf("No access token in response for email: %s\n", email)
		return errs.BadRequest("No access token in response")
	}

	fmt.Printf("Access token: %s\n", verifyResp.Session.AccessToken)

	// Step 2: Update password with the access token
	passwordPayload := map[string]string{
		"password": newPassword,
	}

	fmt.Printf("Password payload: %v\n", passwordPayload)

	passwordBytes, _ := json.Marshal(passwordPayload)
	passwordReq, _ := http.NewRequest("PUT", fmt.Sprintf("%s/auth/v1/user", supabaseURL), bytes.NewBuffer(passwordBytes))
	passwordReq.Header.Set("Content-Type", "application/json")
	passwordReq.Header.Set("apikey", apiKey)
	passwordReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err = Client.Do(passwordReq)
	if err != nil {
		fmt.Printf("Failed to update password request: %v\n", err)
		return errs.BadRequest("Failed to update password")
	}
	defer res.Body.Close()

	body, _ = io.ReadAll(res.Body)
	fmt.Printf("Update password response body: %s\n", string(body))
	if res.StatusCode != http.StatusOK {
		fmt.Printf("Failed to update password: %d, %s\n", res.StatusCode, string(body))
		return errs.BadRequest(fmt.Sprintf("Failed to update password: %d, %s", res.StatusCode, string(body)))
	}

	fmt.Printf("Password update successful for email: %s\n", email)

	return nil
}
