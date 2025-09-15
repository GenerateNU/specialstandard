package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	StartDateTime time.Time  `json:"start_datetime" db:"start_datetime"`
	EndDateTime   time.Time  `json:"end_datetime" db:"end_datetime"`
	TherapistID   uuid.UUID  `json:"therapist_id" db:"therapist_id"`
	Notes         *string    `json:"notes" db:"notes"`
	CreatedAt     *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" db:"updated_at"`
}

type PostSessionInput struct {
	StartTime   time.Time `json:"start_datetime" validate:"required"`
	EndTime     time.Time `json:"end_datetime" validate:"required"`
	TherapistID uuid.UUID `json:"therapist_id" validate:"required"`
	Notes       *string   `json:"notes"`
}

type PatchSessionInput struct {
	StartTime   *time.Time `json:"start_datetime"`
	EndTime     *time.Time `json:"end_datetime"`
	TherapistID *uuid.UUID `json:"therapist_id"`
	Notes       *string    `json:"notes"`
}

// SessionQueryParams for filtering sessions (for future use)
type SessionQueryParams struct {
	TherapistID *uuid.UUID `query:"therapist_id"`
	FromDate    *time.Time `query:"from_date"`
	ToDate      *time.Time `query:"to_date"`
	Limit       *int       `query:"limit"`
	Offset      *int       `query:"offset"`
}
