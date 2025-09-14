package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	GetSessions(ctx context.Context) ([]models.Session, error)
}

type ThemeRepository interface {
	CreateTheme(ctx context.Context, theme *models.CreateThemeInput) (*models.Theme, error)
}

type TherapistRepository interface {
	GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error)
	GetTherapists(ctx context.Context) ([]models.Therapist, error)
	CreateTherapist(ctx context.Context, therapist *models.CreateTherapistInput) (*models.Therapist, error)
	DeleteTherapist(ctx context.Context, therapistID string) (string, error)
	PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error)
}

type Repository struct {
	db        *pgxpool.Pool
	Session   SessionRepository
	Theme     ThemeRepository
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
		db:        db,
		Session:   schema.NewSessionRepository(db),
		Theme:     schema.NewThemeRepository(db),
		Therapist: schema.NewTherapistRepository(db),
	}
}
