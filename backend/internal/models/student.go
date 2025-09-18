package models

import (
	"time"
	"github.com/google/uuid"
)

type Student struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	DOB         *time.Time `json:"dob,omitempty" db:"dob"`               
	TherapistID uuid.UUID  `json:"therapist_id" db:"therapist_id"`
	Grade       *string    `json:"grade,omitempty" db:"grade"`          
	IEP         *string    `json:"iep,omitempty" db:"iep"`              
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// Input struct for creating students
type CreateStudentInput struct {
	FirstName   string  `json:"first_name" validate:"required"`
	LastName    string  `json:"last_name" validate:"required"`
	DOB         *string `json:"dob,omitempty"`                         
	TherapistID string  `json:"therapist_id" validate:"required"`      
	Grade       *string `json:"grade,omitempty"`                       
	IEP         *string `json:"iep,omitempty"`                         
}

// Input struct for updating students (all fields optional for PATCH)
type UpdateStudentInput struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	DOB         *string `json:"dob,omitempty"`
	TherapistID *string `json:"therapist_id,omitempty"`
	Grade       *string `json:"grade,omitempty"`                
	IEP         *string `json:"iep,omitempty"`                       
}