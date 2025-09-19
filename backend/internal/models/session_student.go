package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStudent struct {
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	StudentID uuid.UUID `json:"student_id" db:"student_id"`
	Present   bool      `json:"present" db:"present"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSessionStudentInput struct {
	SessionID uuid.UUID `json:"session_id" validate:"required,uuid"`
	StudentID uuid.UUID `json:"student_id" validate:"required,uuid"`
	Present   bool      `json:"present" validate:"required"`
	Notes     *string   `json:"notes,omitempty"`
}

type PatchSessionStudentInput struct {
	Present *bool   `json:"present,omitempty"`
	Notes   *string `json:"notes,omitempty"`
}
