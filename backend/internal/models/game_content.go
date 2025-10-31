package models

import (
	"time"

	"github.com/google/uuid"
)

type GameContent struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Category  string     `json:"category" db:"category"`
	Level     int        `json:"level" db:"level"`
	Options   []string   `json:"options" db:"options"`
	Answer    string     `json:"answer" db:"answer"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type GetGameContentRequest struct {
	Category string `query:"category" validate:"required,omitempty,oneof=sequencing following_directions wh_questions true_false concepts_sorting"`
	Level    int    `query:"level" validate:"required,omitempty,gte=0,lte=12"`
	Count    int    `query:"count" validate:"required,omitempty,gte=2"`
}
