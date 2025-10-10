package models

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID         uuid.UUID  `json:"id"`
	ThemeID    uuid.UUID  `json:"theme_id"`
	GradeLevel *int       `json:"grade_level"`
	Date       *time.Time `json:"date"`
	Type       *string    `json:"type"`
	Title      *string    `json:"title"`
	Category   *string    `json:"category"`
	Content    *string    `json:"content"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at"`
}

type ResourceWithTheme struct {
	Resource
	Theme ThemeInfo `json:"theme"`
}

type ResourceBody struct {
	ThemeID    uuid.UUID  `json:"theme_id"`
	GradeLevel *int       `json:"grade_level" validate:"omitempty,oneof=0 1 2 3 4 5 6 7 8 9 10 11 12"`
	Date       *time.Time `json:"date"`
	Type       *string    `json:"type"`
	Title      *string    `json:"title"`
	Category   *string    `json:"category"`
	Content    *string    `json:"content"`
}

type UpdateResourceBody struct {
	ThemeID    *uuid.UUID `json:"theme_id"`
	GradeLevel *int       `json:"grade_level" validate:"omitempty,oneof=0 1 2 3 4 5 6 7 8 9 10 11 12"`
	Date       *time.Time `json:"date"`
	Type       *string    `json:"type"`
	Title      *string    `json:"title"`
	Category   *string    `json:"category"`
	Content    *string    `json:"content"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
