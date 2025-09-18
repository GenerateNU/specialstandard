package models

import (
	"time"

	"github.com/google/uuid"
)

type Therapist struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTherapistInput struct {
	FirstName string `json:"first_name" validate:"required,min=1,max=255"`
	LastName  string `json:"last_name" validate:"required,min=1,max=255"`
	Email     string `json:"email" validate:"required,min=1,max=255"`
}

type UpdateTherapist struct {
	FirstName *string `json:"first_name" validate:"omitempty,min=1,max=255"` 
	LastName  *string `json:"last_name" validate:"omitempty,min=1,max=255"`
	Email     *string `json:"email" validate:"omitempty,min=1,max=255"`
	Active    *bool   `json:"active" validate:"omitempty,min=1,max=255"`
}
