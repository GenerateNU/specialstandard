package models

import (
	"specialstandard/internal/utils"
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	FirstName   string     `json:"first_name" db:"first_name"`
	LastName    string     `json:"last_name" db:"last_name"`
	DOB         *time.Time `json:"dob,omitempty" db:"dob"`
	TherapistID uuid.UUID  `json:"therapist_id" db:"therapist_id"`
	SchoolID    int        `json:"school_id" db:"school_id"`
	SchoolName  *string    `json:"school_name,omitempty" db:"school_name"` // GET
	DistrictID  *int       `json:"district_id,omitempty" db:"district_id"` // GET requests ONLY
	Grade       *int       `json:"grade,omitempty" db:"grade"`
	IEP         *string    `json:"iep,omitempty" db:"iep"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type CreateStudentInput struct {
	FirstName   string  `json:"first_name" validate:"required,min=1,max=100"`
	LastName    string  `json:"last_name" validate:"required,min=1,max=100"`
	DOB         *string `json:"dob,omitempty" validate:"omitempty,datetime=2006-01-02"`
	TherapistID string  `json:"therapist_id" validate:"required,uuid"`
	Grade       *int    `json:"grade,omitempty" validate:"omitempty,oneof=0 1 2 3 4 5 6 7 8 9 10 11 12"`
	IEP         *string `json:"iep,omitempty"`
	SchoolID    int     `json:"school_id" validate:"required,min=1"`
}

type GetStudentsQuery struct {
	Grade       *int   `query:"grade" validate:"omitempty,oneof=-1 0 1 2 3 4 5 6 7 8 9 10 11 12"`
	TherapistID string `query:"therapist_id" validate:"omitempty,uuid"`
	Name        string `query:"name" validate:"omitempty"`
	utils.Pagination
}

type UpdateStudentInput struct {
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	DOB         *string `json:"dob,omitempty"`
	TherapistID *string `json:"therapist_id,omitempty"`
	SchoolID    *int    `json:"school_id,omitempty"`
	Grade       *int    `json:"grade,omitempty" validate:"omitempty,oneof=-1 0 1 2 3 4 5 6 7 8 9 10 11 12"`
	IEP         *string `json:"iep,omitempty"`
}

type PromoteStudentsInput struct {
	TherapistID        uuid.UUID   `json:"therapist_id" validate:"required"`
	ExcludedStudentIDs []uuid.UUID `json:"excluded_student_ids" validate:"dive"`
}
