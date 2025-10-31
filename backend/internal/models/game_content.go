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
