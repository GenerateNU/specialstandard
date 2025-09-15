package schema

import (
	"context"
	"database/sql"
	"errors"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func (r *SessionRepository) GetSessions(ctx context.Context) ([]models.Session, error) {
	query := `
	SELECT id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at
	FROM session`

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

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	query := `
	SELECT id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at
	FROM session
	WHERE id = $1`

	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	// Using CollectExactlyOneRow because we expect exactly one session with this ID
	session, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.Session])
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionRepository) DeleteSessions(ctx context.Context, id uuid.UUID) (string, error) {
	query := `DELETE FROM session WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return "", err
	}

	// If No Rows are deleted/affected.
	if cmdTag.RowsAffected() == 0 {
		return "not found", nil
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
		&session.StartDateTime,
		&session.EndDateTime,
		&session.TherapistID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return session, nil
}

func (r *SessionRepository) PatchSessions(ctx context.Context, id uuid.UUID, input *models.PatchSessionInput) (*models.Session, error) {
	session := &models.Session{}

	query := `UPDATE session
				SET
					start_datetime = COALESCE($1, start_datetime),
					end_datetime = COALESCE($2, end_datetime),
					therapist_id = COALESCE($3, therapist_id),
					notes = COALESCE($4, notes)
				WHERE id = $5
				RETURNING id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.StartTime, input.EndTime, input.TherapistID, input.Notes, id)

	if err := row.Scan(
		&session.ID,
		&session.StartDateTime,
		&session.EndDateTime,
		&session.TherapistID,
		&session.Notes,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.NotFound("Session Not Found")
		}
		return nil, err
	}

	return session, nil
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}
