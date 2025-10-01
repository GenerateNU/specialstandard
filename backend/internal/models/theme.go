package models

import (
	"time"

	"github.com/google/uuid"
)

type Theme struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"theme_name"`
	Month     int        `json:"month" db:"month"`
	Year      int        `json:"year" db:"year"`
	CreatedAt *time.Time `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type CreateThemeInput struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Month int    `json:"month" validate:"required,gte=1,lte=12"`
	Year  int    `json:"year" validate:"required,gte=1900,lte=2100"`
}

type UpdateThemeInput struct {
	Name  *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Month *int    `json:"month,omitempty" validate:"omitempty,gte=1,lte=12"`
	Year  *int    `json:"year,omitempty" validate:"omitempty,gte=1900,lte=2100"`
}

type ThemeFilter struct {
	Month  *int    `query:"month" validate:"omitempty,gte=1,lte=12"`
	Year   *int    `query:"year" validate:"omitempty,gte=1900,lte=2100"`
	Search *string `query:"search" validate:"omitempty,max=255"`
}
