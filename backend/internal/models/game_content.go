package models

import (
	"time"

	"github.com/aws/smithy-go/ptr"
	"github.com/google/uuid"
)

type GameContent struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	ThemeID         uuid.UUID  `json:"theme_id" db:"theme_id"`
	Week            int        `json:"week" db:"week"`
	Category        *string    `json:"category" db:"category"`
	QuestionType    string     `json:"question_type" db:"question_type"`
	DifficultyLevel int        `json:"difficulty_level" db:"difficulty_level"`
	Question        string     `json:"question" db:"question"`
	Options         []string   `json:"options" db:"options"`
	Answer          string     `json:"answer" db:"answer"`
	CreatedAt       *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at" db:"updated_at"`
}

type GetGameContentRequest struct {
	ThemeID         *uuid.UUID `query:"theme_id" validate:"omitempty,uuid"`
	Category        *string    `query:"category" validate:"omitempty,oneof=receptive_language expressive_language social_pragmatic_language speech"`
	QuestionType    *string    `query:"question_type" validate:"omitempty,oneof=sequencing following_directions wh_questions true_false concepts_sorting fill_in_the_blank categorical_language emotions teamwork_talk express_excitement_interest fluency articulation_s articulation_l"`
	DifficultyLevel *int       `query:"difficulty_level" validate:"omitempty,gte=1"`
	Count           *int       `query:"count" validate:"omitempty,gte=2"`
}

const (
	defaultCount int = 4
)

func NewGetGameContentRequest() GetGameContentRequest {
	return GetGameContentRequest{
		Count: ptr.Int(defaultCount),
	}
}
