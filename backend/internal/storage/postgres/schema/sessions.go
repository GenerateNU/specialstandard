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

func (r *SessionRepository) DeleteSessions(ctx context.Context, id int) (string, error) {
	session := &models.Session{}

	query := `DELETE FROM sessions WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	if err := row.Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.TherapistID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return "", err
	}

	return "Deleted the Session Successfully!", nil
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}
