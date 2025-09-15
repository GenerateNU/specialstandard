package models

import (
	"time"
	"github.com/google/uuid"
)


type Student struct {
	ID          uuid.UUID `json:"id" db:"id"`
	FirstName   string    `json:"first_name" db:"first_name"`
	LastName    string    `json:"last_name" db:"last_name"`
	DOB         time.Time `json:"dob" db:"dob"`
	TherapistID uuid.UUID `json:"therapist_id" db:"therapist_id"`
	Grade       string    `json:"grade" db:"grade"`
	IEP         string    `json:"iep" db:"iep"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}