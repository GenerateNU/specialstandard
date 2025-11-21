package models

import (
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/google/uuid"
)

type GameContent struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	ThemeID             uuid.UUID  `json:"theme_id" db:"theme_id"`
	Week                int        `json:"week" db:"week"`
	Category            *string    `json:"category" db:"category"`
	QuestionType        string     `json:"question_type" db:"question_type"`
	DifficultyLevel     int        `json:"difficulty_level" db:"difficulty_level"`
	Question            string     `json:"question" db:"question"`
	Options             []string   `json:"options" db:"options"`
	Answer              string     `json:"answer" db:"answer"`
	CreatedAt           *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at" db:"updated_at"`
	ExerciseType        string     `json:"exercise_type" db:"exercise_type"`
	ApplicableGameTypes []string   `json:"applicable_game_types" db:"applicable_game_types"`
}

type GetGameContentRequest struct {
	ThemeID             *uuid.UUID `query:"theme_id" validate:"omitempty,uuid"`
	Category            *string    `query:"category" validate:"omitempty,oneof=receptive_language expressive_language social_pragmatic_language speech"`
	QuestionType        *string    `query:"question_type" validate:"omitempty,oneof=sequencing following_directions wh_questions true_false concepts_sorting fill_in_the_blank categorical_language emotions teamwork_talk express_excitement_interest fluency articulation_s articulation_l"`
	DifficultyLevel     *int       `query:"difficulty_level" validate:"omitempty,gte=1"`
	QuestionCount       *int       `query:"question_count" validate:"omitempty,gte=2"`
	WordsCount          *int       `query:"words_count" validate:"omitempty,gte=2"`
	ExerciseType        *string    `query:"exercise_type" validate:"omitempty,oneof=game pdf"`
	ApplicableGameTypes *[]string  `query:"applicable_game_types" validate:"omitempty,dive,oneof=sequencing following_directions wh_questions true_false concepts_sorting fill_in_the_blank categorical_language emotions teamwork_talk express_excitement_interest fluency articulation_s articulation_l"`
}

const (
	defaultQuestionCount int = 5
	defaultWordsCount    int = 4
)

func NewGetGameContentRequest() GetGameContentRequest {
	return GetGameContentRequest{
		QuestionCount: ptr.Int(defaultQuestionCount),
		WordsCount:    ptr.Int(defaultWordsCount),
	}
}
