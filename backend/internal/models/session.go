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

type Repetition struct {
	RecurStart  time.Time `json:"recur_start" validate:"required"`
	RecurEnd    time.Time `json:"recur_end" validate:"required,gtfield=RecurStart"`
	EveryNWeeks int       `json:"every_n_weeks" validate:"required,gte=1"`
}

type PostSessionInput struct {
	StartTime   time.Time    `json:"start_datetime" validate:"required"`
	EndTime     time.Time    `json:"end_datetime" validate:"required"`
	TherapistID uuid.UUID    `json:"therapist_id" validate:"required"`
	Notes       *string      `json:"notes"`
	Repetition  *Repetition  `json:"repetition" validate:"omitempty"`
	StudentIDs  *[]uuid.UUID `json:"student_ids" validate:"omitempty,dive,uuid"`
}

type PatchSessionInput struct {
	StartTime   *time.Time `json:"start_datetime"`
	EndTime     *time.Time `json:"end_datetime"`
	TherapistID *uuid.UUID `json:"therapist_id"`
	Notes       *string    `json:"notes"`
}

type GetSessionRequest struct {
	StartTime  *time.Time `query:"startdate" validate:"omitempty"`
	EndTime    *time.Time `query:"enddate" validate:"omitempty"`
	Month      *int       `query:"month" validate:"omitempty,gte=1,lte=12"`
	Year       *int       `query:"year" validate:"omitempty,gte=1776,lte=2200"`
	StudentIDs *[]string  `query:"id" validate:"omitempty"`
}

// This is what repository uses
type GetSessionRepositoryRequest struct {
	StartTime  *time.Time   `validate:"omitempty"`
	EndTime    *time.Time   `validate:"omitempty"`
	Month      *int         `validate:"omitempty,gte=1,lte=12"`
	Year       *int         `validate:"omitempty,gte=1776,lte=2200"`
	StudentIDs *[]uuid.UUID `validate:"omitempty"`
}

type GetStudentSessionsRequest struct {
	Month   *int  `query:"month" validate:"omitempty,gte=1,lte=12"`
	Year    *int  `query:"year" validate:"omitempty,gte=1776,lte=2200"`
	Present *bool `query:"present" validate:"omitempty"`
}

// This is what repository uses
type GetStudentSessionsRepositoryRequest struct {
	StartDate *time.Time `validate:"omitempty"`
	EndDate   *time.Time `validate:"omitempty"`
	Month     *int       `validate:"omitempty,gte=1,lte=12"`
	Year      *int       `validate:"omitempty,gte=1776,lte=2200"`
	Present   *bool      `validate:"omitempty"`
}
