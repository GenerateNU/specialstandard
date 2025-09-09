package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID  `json:"id"`
	StartTime   string     `json:"start_time"`
	EndTime     string     `json:"end_time"`
	TherapistID uuid.UUID  `json:"therapist_id"`
	Notes       *string    `json:"notes"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type PostSessionInput struct {
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	TherapistID uuid.UUID `json:"therapist_id"`
	Notes       *string   `json:"notes"`
}

type PatchSessionInput struct {
	StartTime   *string    `json:"start_time"`
	EndTime     *string    `json:"end_time"`
	TherapistID *uuid.UUID `json:"therapist_id"`
	Notes       *string    `json:"notes"`
}
