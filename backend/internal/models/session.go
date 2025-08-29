package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID  `json:"id"`
	TherapistID uuid.UUID  `json:"therapist_id"`
	SessionDate time.Time  `json:"session_date"`
	StartTime   *string    `json:"start_time"`
	EndTime     *string    `json:"end_time"`
	Notes       *string    `json:"notes"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
