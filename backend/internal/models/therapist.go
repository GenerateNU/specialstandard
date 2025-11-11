package models

import (
	"time"

	"github.com/google/uuid"
)

type Therapist struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Email        string    `json:"email"`
	Active       bool      `json:"active"`
	Schools      []int     `json:"schools" db:"schools"`
	DistrictID   *int      `json:"district_id"`
	SchoolNames  *[]string `json:"school_names,omitempty" db:"-"`
	DistrictName *string   `json:"district_name,omitempty" db:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateTherapistInput struct {
	ID         uuid.UUID `json:"id" validate:"required"`
	FirstName  string    `json:"first_name" validate:"required,min=1,max=255"`
	LastName   string    `json:"last_name" validate:"required,min=1,max=255"`
	Schools    []int     `json:"schools" validate:"required,dive,min=1"`
	DistrictID *int      `json:"district_id" validate:"omitempty,min=1"`
	Email      string    `json:"email" validate:"required,min=1,max=255"`
}

type UpdateTherapist struct {
	FirstName  *string `json:"first_name" validate:"omitempty,min=1,max=255"`
	LastName   *string `json:"last_name" validate:"omitempty,min=1,max=255"`
	Schools    *[]int  `json:"schools" validate:"omitempty,dive,min=1"`
	DistrictID *int    `json:"district_id" validate:"omitempty,min=1"`
	Email      *string `json:"email" validate:"omitempty,min=1,max=255"`
	Active     *bool   `json:"active" validate:"omitempty"`
}
