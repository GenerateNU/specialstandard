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

	query := `DELETE FROM session WHERE id = $1`
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

func (r *SessionRepository) PostSessions(ctx context.Context, input *models.PostSessionInput) (*models.Session, error) {
	session := &models.Session{}

	query := `INSERT INTO session (start_datetime, end_datetime, therapist_id, notes)
				VALUES ($1, $2, $3, $4)
				RETURNING id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.StartTime, input.EndTime, input.TherapistID, input.Notes)
	// Scan into Session model to return the one we just inserted.
	if err := row.Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.TherapistID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) PatchSessions(ctx context.Context, id int, input *models.PatchSessionInput) (*models.Session, error) {
	session := &models.Session{}

	query := `UPDATE session
				SET
					start_datetime = COALESCE($1, start_datetime),
					end_datetime = COALESCE($2, end_datetime),
					therapist_id = COALESCE($3, therapist_id),
					notes = COALESCE($4, notes)
				WHERE notes = $5
				RETURNING id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.StartTime, input.EndTime, input.TherapistID, input.Notes, id)

	if err := row.Scan(
		&session.ID,
		&session.StartTime,
		&session.EndTime,
		&session.TherapistID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return session, nil
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}
