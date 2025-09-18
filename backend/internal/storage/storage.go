package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
)

type SessionRepository interface {
	GetSessions(ctx context.Context) ([]models.Session, error)
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
}

type Repository struct {
	db      *pgxpool.Pool
	Session SessionRepository
	Student StudentRepository
	Theme   ThemeRepository
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
		Theme:   schema.NewThemeRepository(db),
	}
}
