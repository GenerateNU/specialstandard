package models

import (
	"time"
)

// VerificationCode represents a verification code in the database
type VerificationCode struct {
	ID        string        `db:"id"`
	UserID    string     `db:"user_id"`
	Code      string     `db:"code"`
	ExpiresAt time.Time  `db:"expires_at"`
	CreatedAt time.Time  `db:"created_at"`
	Used      bool 		 `db:"used"`
}

type SendCodeResponse struct {
	Success   bool   `json:"success"`
	MessageID string `json:"messageId,omitempty"`
	Error     string `json:"error,omitempty"`
}

type VerifyCodeRequest struct {
	Code string `json:"code"`
}

type VerifyCodeResponse struct {
	Success  bool   `json:"success"`
	Verified bool   `json:"verified,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}