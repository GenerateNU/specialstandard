package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	GetSessions(ctx context.Context) ([]models.Session, error)
	GetSessionByID(ctx context.Context, id string) (*models.Session, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	PostSession(ctx context.Context, session *models.PostSessionInput) (*models.Session, error)
	PatchSession(ctx context.Context, id uuid.UUID, session *models.PatchSessionInput) (*models.Session, error)
}

type StudentRepository interface {
	GetStudents(ctx context.Context) ([]models.Student, error)
	GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error)
	AddStudent(ctx context.Context, student models.Student) (models.Student, error)
	UpdateStudent(ctx context.Context, student models.Student) (models.Student, error)
	DeleteStudent(ctx context.Context, id uuid.UUID) error
}
  
type ThemeRepository interface {
	CreateTheme(ctx context.Context, theme *models.CreateThemeInput) (*models.Theme, error)
	GetThemes(ctx context.Context) ([]models.Theme, error)
	GetThemeByID(ctx context.Context, id uuid.UUID) (*models.Theme, error)
	UpdateTheme(ctx context.Context, id uuid.UUID, theme *models.UpdateThemeInput) (*models.Theme, error)
	DeleteTheme(ctx context.Context, id uuid.UUID) error
}

type TherapistRepository interface {
	GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error)
	GetTherapists(ctx context.Context) ([]models.Therapist, error)
	CreateTherapist(ctx context.Context, therapist *models.CreateTherapistInput) (*models.Therapist, error)
	DeleteTherapist(ctx context.Context, therapistID string) error
	PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error)
}

type Repository struct {
	db      *pgxpool.Pool
	Session SessionRepository
	Student StudentRepository
	Theme   ThemeRepository
	Therapist TherapistRepository
}

func (r *Repository) Close() error {
	r.db.Close()
	return nil
}

func (r *Repository) GetDB() *pgxpool.Pool {
	return r.db
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db:      db,
		Session: schema.NewSessionRepository(db),
		Student: schema.NewStudentRepository(db),
		Theme:     schema.NewThemeRepository(db),
		Therapist: schema.NewTherapistRepository(db),
	}
}
