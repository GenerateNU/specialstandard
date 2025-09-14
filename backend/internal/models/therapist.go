package models

import (
	"time"

	"github.com/google/uuid"
)

type Therapist struct {
	ID         uuid.UUID `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Email      string    `json:"email"`
	Active     bool      `json:"active"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

type CreateTherapistInput struct {
	First_name string `json:"first_name"`
	Last_name  string `json:"last_name"`
	Email      string `json:"email"`
}

type UpdateTherapist struct {
	First_name *string `json:"first_name"`
	Last_name  *string `json:"last_name"`
	Email      *string `json:"email"`
	Active     *bool   `json:"active"`
}
