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
	SessionID       *uuid.UUID `query:"session_id" validate:"omitempty,uuid"`
	StudentID       *uuid.UUID `query:"student_id" validate:"omitempty,uuid"`
	Category        *string    `query:"category" validate:"omitempty,oneof=receptive_language expressive_language social_pragmatic_language speech"`
	QuestionType    *string    `query:"question_type" validate:"omitempty,oneof=sequencing following_directions wh_questions true_false concepts_sorting fill_in_the_blank categorical_language emotions teamwork_talk express_excitement_interest fluency articulation_s articulation_l"`
	DifficultyLevel *int       `query:"difficulty_level" validate:"omitempty,gte=1"`
	ExerciseType    *string    `query:"exercise_type" validate:"omitempty,oneof=game pdf"`
	GameType        *string    `query:"game_type" validate:"omitempty,dive"`
	DateFrom        *time.Time `query:"date_from" validate:"omitempty"`
	DateTo          *time.Time `query:"date_to" validate:"omitempty"`
}

type PostGameResult struct {
	SessionStudentID       int       `json:"session_student_id" validate:"gte=0"` // Remove validate:"required"
	ContentID              uuid.UUID `json:"content_id" validate:"required,uuid"`
	TimeTakenSec           int       `json:"time_taken_sec" validate:"gte=0"` // Remove required
	Completed              *bool     `json:"completed,omitempty"`
	CountIncorrectAttempts int       `json:"count_of_incorrect_attempts" validate:"gte=0"` // Remove required
	IncorrectAttempts      *[]string `json:"incorrect_attempts,omitempty" validate:"omitempty,dive"`
}
