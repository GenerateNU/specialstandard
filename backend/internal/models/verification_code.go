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