package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepositoryImpl struct {
	db *pgxpool.Pool
}

func (r *SessionRepositoryImpl) GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest) ([]models.Session, error) {
	query := `
	SELECT id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at
	FROM session`

	conditions := []string{}
	args := []interface{}{}
	argCount := 1

	if filter != nil {
		if filter.Month != nil && filter.Year != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(MONTH FROM start_datetime) = $%d AND EXTRACT(YEAR FROM start_datetime) = $%d", argCount, argCount+1))
			args = append(args, *filter.Month, *filter.Year)
			argCount += 2
		} else if filter.Year != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM start_datetime) = $%d", argCount))
			args = append(args, *filter.Year)
			argCount++
		} else if filter.Month != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(MONTH FROM start_datetime) = $%d", argCount))
			args = append(args, *filter.Month)
			argCount++
		}

		if filter.StartTime != nil {
			conditions = append(conditions, fmt.Sprintf("start_datetime >= $%d", argCount))
			args = append(args, *filter.StartTime)
			argCount++
		}

		if filter.EndTime != nil {
			conditions = append(conditions, fmt.Sprintf("end_datetime <= $%d", argCount))
			args = append(args, *filter.EndTime)
			argCount++
		}

		if filter.StudentIDs != nil && len(*filter.StudentIDs) > 0 {
			// For each student ID, add a condition that checks if that specific student exists in the session
			for _, studentID := range *filter.StudentIDs {
				conditions = append(conditions, fmt.Sprintf(
					"EXISTS (SELECT 1 FROM session_student WHERE session_id = session.id AND student_id = $%d)",
					argCount,
				))
				args = append(args, studentID)
				argCount++
			}
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(` ORDER BY start_datetime DESC LIMIT $%d OFFSET $%d`, argCount, argCount+1)

	args = append(args, pagination.Limit, pagination.GettOffset())

	rows, err := r.db.Query(ctx, query, args...)
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

func (r *SessionRepositoryImpl) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
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

func (r *SessionRepositoryImpl) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM session WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SessionRepositoryImpl) PostSession(ctx context.Context, input *models.PostSessionInput) (*models.Session, error) {
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

func (r *SessionRepositoryImpl) PatchSession(ctx context.Context, id uuid.UUID, input *models.PatchSessionInput) (*models.Session, error) {
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
		return nil, err
	}

	return session, nil
}

func NewSessionRepository(db *pgxpool.Pool) *SessionRepositoryImpl {
	return &SessionRepositoryImpl{
		db,
	}
}

func (r *SessionRepositoryImpl) GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.SessionStudentsOutput, error) {
	query := `
	SELECT ss.session_id, ss.present, ss.notes, ss.created_at, ss.updated_at,
	       s.id, s.first_name, s.last_name, s.dob, s.therapist_id, 
	       s.grade, s.iep, s.created_at, s.updated_at
	FROM session_student ss
	JOIN student s ON ss.student_id = s.id
	WHERE ss.session_id = $1
	LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, sessionID, pagination.Limit, pagination.GettOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessionStudents []models.SessionStudentsOutput
	for rows.Next() {
		var result models.SessionStudentsOutput
		var student models.Student

		err := rows.Scan(
			&result.SessionID, &result.Present, &result.Notes, &result.CreatedAt, &result.UpdatedAt,
			&student.ID, &student.FirstName, &student.LastName, &student.DOB, &student.TherapistID,
			&student.Grade, &student.IEP, &student.CreatedAt, &student.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Student = student
		sessionStudents = append(sessionStudents, result)
	}

	return sessionStudents, nil
}
