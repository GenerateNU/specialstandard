package models

import (
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Month     int        `json:"month"`
	Year      int        `json:"year"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type CreateThemeInput struct {
	Name  string `json:"name"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
}
