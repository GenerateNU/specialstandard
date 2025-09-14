package storage

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/postgres/schema"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	GetSessions(ctx context.Context) ([]models.Session, error)
	DeleteSessions(ctx context.Context, id int) (string, error)
	PostSessions(ctx context.Context, session *models.PostSessionInput) (*models.Session, error)
	PatchSessions(ctx context.Context, id int, session *models.PatchSessionInput) (*models.Session, error)
}

type ThemeRepository interface {
	CreateTheme(ctx context.Context, theme *models.CreateThemeInput) (*models.Theme, error)
}

type Repository struct {
	db      *pgxpool.Pool
	Session SessionRepository
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
		Theme:   schema.NewThemeRepository(db),
	}
}
