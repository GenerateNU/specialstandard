package models

import (
	"time"

	"github.com/google/uuid"
)

type SessionStudent struct {
	ID        int       `json:"id" db:"id"`
	SessionID uuid.UUID `json:"session_id" db:"session_id"`
	StudentID uuid.UUID `json:"student_id" db:"student_id"`
	Present   bool      `json:"present" db:"present"`
	Notes     *string   `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type SessionStudentsOutput struct {
	Student   Student         `json:"student" db:"student"`
	SessionID uuid.UUID       `json:"session_id" db:"session_id"`
	Present   bool            `json:"present" db:"present"`
	Notes     *string         `json:"notes,omitempty" db:"notes"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt time.Time       `json:"updated_at" db:"updated_at"`
	Ratings   []SessionRating `json:"ratings" db:"ratings"`
}

type StudentSessionsOutput struct {
	Session   Session   `json:"session"`
	StudentID uuid.UUID `json:"student_id"`
	Present   bool      `json:"present"`
	Notes     *string   `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StudentSessionsWithRatingsOutput struct {
	SessionID   uuid.UUID       `json:"session_id"`
	StudentID   uuid.UUID       `json:"student_id"`
	SessionDate time.Time       `json:"session_date"`
	Ratings     []SessionRating `json:"ratings"`
}

type CreateSessionStudentInput struct {
	SessionIDs []uuid.UUID `json:"session_ids" validate:"required,min=1,dive,uuid"`
	StudentIDs []uuid.UUID `json:"student_ids" validate:"required,min=1,dive,uuid"`
	Present    bool        `json:"present"`
	Notes      *string     `json:"notes,omitempty"`
}

type PatchSessionStudentInput struct {
	SessionID uuid.UUID    `json:"session_id" validate:"required,uuid"`
	StudentID uuid.UUID    `json:"student_id" validate:"required,uuid"`
	Present   *bool        `json:"present,omitempty"`
	Notes     *string      `json:"notes,omitempty"`
	Ratings   *[]RateInput `json:"ratings" validate:"required,dive"`
}

type DeleteSessionStudentInput struct {
	SessionID uuid.UUID `json:"session_id" validate:"required,uuid"`
	StudentID uuid.UUID `json:"student_id" validate:"required,uuid"`
}

type RateInput struct {
	Category    string `json:"category" validate:"required,oneof=visual_cue verbal_cue gestural_cue engagement"`
	Level       string `json:"level" validate:"required,oneof=minimal moderate maximal low high"`
	Description string `json:"description" validate:"required"`
}

type SessionRating struct {
	Category    *string `json:"category" validate:"oneof=visual_cue verbal_cue gestural_cue engagement"`
	Level       *string `json:"level" validate:"oneof=minimal moderate maximal low high"`
	Description *string `json:"description"`
}

type PatchSessionStudentRatingsOutput struct {
	SessionID uuid.UUID       `json:"session_id"`
	StudentID uuid.UUID       `json:"student_id"`
	Present   *bool           `json:"present,omitempty"`
	Notes     *string         `json:"notes,omitempty"`
	Ratings   []SessionRating `json:"ratings"`
}
