package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func (r *SessionRepository) GetSessions(ctx context.Context) ([]models.Session, error) {
	query := `
	SELECT id, therapist_id, session_date, start_time, end_time, notes, created_at, updated_at
	FROM sessions`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Session])
	if err != nil {
		return nil, err
	}
	return sessions, nil
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}
