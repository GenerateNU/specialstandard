package models

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID         uuid.UUID  `json:"id"`
	ThemeID    uuid.UUID  `json:"theme_id"`
	GradeLevel *string    `json:"grade_level"`
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
	GradeLevel *string    `json:"grade_level"`
	Date       *time.Time `json:"date"`
	Type       *string    `json:"type"`
	Title      *string    `json:"title"`
	Category   *string    `json:"category"`
	Content    *string    `json:"content"`
}

type UpdateResourceBody struct {
	ThemeID    *uuid.UUID `json:"theme_id"`
	GradeLevel *string    `json:"grade_level"`
	Date       *time.Time `json:"date"`
	Type       *string    `json:"type"`
	Title      *string    `json:"title"`
	Category   *string    `json:"category"`
	Content    *string    `json:"content"`
	UpdatedAt  *time.Time `json:"updated_at"`
}
