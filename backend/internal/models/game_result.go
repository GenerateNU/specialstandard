package models

import (
	"time"

	"github.com/google/uuid"
)

type GameResult struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	SessionID      uuid.UUID  `json:"session_id" db:"session_id"`
	StudentID      uuid.UUID  `json:"student_id" db:"student_id"`
	ContentID      uuid.UUID  `json:"content_id" db:"content_id"`
	TimeTaken      int        `json:"time_taken" db:"time_taken"`
	Completed      bool       `json:"completed" db:"completed"`
	IncorrectTries int        `json:"incorrect_tries" db:"incorrect_tries"`
	CreatedAt      *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at" db:"updated_at"`
}

type PostGameResult struct {
	SessionID      uuid.UUID `json:"session_id" validate:"required,uuid"`
	StudentID      uuid.UUID `json:"student_id" validate:"required,uuid"`
	ContentID      uuid.UUID `json:"content_id" validate:"required,uuid"`
	TimeTaken      int       `json:"time_taken" validate:"required,gte=0"`
	Completed      bool      `json:"completed" validate:"required"`
	IncorrectTries int       `json:"incorrect_tries" validate:"required,gte=0"`
}
