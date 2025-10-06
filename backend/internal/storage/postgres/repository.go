package postgres

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/google/uuid"
)

type StudentRepository interface {
	GetStudents(ctx context.Context, grade *int, therapistID uuid.UUID, name string, pagination utils.Pagination) ([]models.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error)
	AddStudent(ctx context.Context, student models.Student) (models.Student, error)
	UpdateStudent(ctx context.Context, student models.Student) (models.Student, error)
	GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination) ([]models.StudentSessionsOutput, error)
	DeleteStudent(ctx context.Context, id uuid.UUID) error
}

type TherapistRepository interface {
	GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error)
	GetTherapists(ctx context.Context, pagination utils.Pagination) ([]models.Therapist, error)
	CreateTherapist(ctx context.Context, input *models.CreateTherapistInput) (*models.Therapist, error)
	DeleteTherapist(ctx context.Context, therapistID string) error
	PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error)
}
