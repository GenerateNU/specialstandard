package models

import (
	"time"

	"github.com/google/uuid"
)

type GameResult struct {
	ID                     uuid.UUID  `json:"id" db:"id"`
	SessionStudentID       int        `json:"session_student_id" db:"session_student_id"`
	ContentID              uuid.UUID  `json:"content_id" db:"content_id"`
	TimeTakenSec           int        `json:"time_taken_sec" db:"time_taken_sec"`
	Completed              bool       `json:"completed" db:"completed"`
	CountIncorrectAttempts int        `json:"count_of_incorrect_attempts" db:"count_of_incorrect_attempts"`
	IncorrectAttempts      *[]string  `json:"incorrect_attempts" db:"incorrect_attempts"`
	CreatedAt              *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              *time.Time `json:"updated_at" db:"updated_at"`
}

type GetGameResultQuery struct {
	SessionID *uuid.UUID `query:"session_id" validate:"omitempty,uuid"`
	StudentID *uuid.UUID `query:"student_id" validate:"omitempty,uuid"`
}

type PostGameResult struct {
	SessionStudentID       int       `json:"session_student_id" validate:"required"`
	ContentID              uuid.UUID `json:"content_id" validate:"required,uuid"`
	TimeTakenSec           int       `json:"time_taken_sec" validate:"required,gte=0"`
	Completed              *bool     `json:"completed,omitempty"`
	CountIncorrectAttempts int       `json:"count_of_incorrect_attempts" validate:"required,gte=0"`
	IncorrectAttempts      *[]string `json:"incorrect_attempts,omitempty" validate:"omitempty,dive"`
}
