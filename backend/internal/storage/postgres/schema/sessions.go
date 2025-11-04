package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/dbinterface"
	"specialstandard/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	db *pgxpool.Pool
}

func (r *SessionRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func (r *SessionRepository) GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest) ([]models.Session, error) {
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

	query += fmt.Sprintf(` ORDER BY start_datetime ASC LIMIT $%d OFFSET $%d`, argCount, argCount+1)

	args = append(args, pagination.Limit, pagination.GetOffset())

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

func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM session WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func addNWeeks(timestamp time.Time, nWeeks int) time.Time {
	return timestamp.Add(time.Duration(24*7*nWeeks) * time.Hour)
}

func (r *SessionRepository) PostSession(ctx context.Context, q dbinterface.Queryable, input *models.PostSessionInput) (*[]models.Session, error) {
	query := `INSERT INTO session (start_datetime, end_datetime, therapist_id, notes)
              VALUES ($1, $2, $3, $4)`
	args := []interface{}{}
	args = append(args, input.StartTime, input.EndTime, input.TherapistID, input.Notes)
	argCount := 5

	if input.Repetition != nil {
		rp := input.Repetition
		if rp.RecurEnd.Before(rp.RecurStart) {
			return nil, fmt.Errorf("invalid repetition range: recur_end (%s) is before recur_start (%s)",
				rp.RecurEnd.Format(time.RFC3339), rp.RecurStart.Format(time.RFC3339))
		}

		startTime := addNWeeks(input.StartTime, rp.EveryNWeeks)
		endTime := addNWeeks(input.EndTime, rp.EveryNWeeks)

		for startTime.Before(rp.RecurEnd) {
			query += fmt.Sprintf(`, ($%d, $%d, $%d, $%d)`, argCount, argCount+1, argCount+2, argCount+3)
			argCount += 4
			args = append(args, startTime, endTime, input.TherapistID, input.Notes)

			startTime = addNWeeks(startTime, rp.EveryNWeeks)
			endTime = addNWeeks(endTime, rp.EveryNWeeks)
		}
	}

	query += ` RETURNING id, start_datetime, end_datetime, therapist_id, notes, created_at, updated_at`

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Session])
	if err != nil {
		return nil, err
	}
	return &sessions, nil
}

func (r *SessionRepository) PatchSession(ctx context.Context, id uuid.UUID, input *models.PatchSessionInput) (*models.Session, error) {
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

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}

func (r *SessionRepository) GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.SessionStudentsOutput, error) {
	query := `
    SELECT ss.session_id, ss.present, ss.notes, ss.created_at, ss.updated_at,
           s.id, s.first_name, s.last_name, s.dob, s.therapist_id, 
           s.grade, s.iep, s.created_at, s.updated_at, 
           sr.level, sr.category, sr.description
    FROM session_student ss
    JOIN student s ON ss.student_id = s.id
    LEFT JOIN session_rating sr ON ss.id = sr.session_student_id
		JOIN session se ON ss.session_id = se.id
    WHERE ss.session_id = $1
    AND s.grade != -1
    ORDER BY se.start_datetime ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, sessionID, pagination.Limit, pagination.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionStudentsMap := make(map[string]*models.SessionStudentsOutput)
	for rows.Next() {
		var result models.SessionStudentsOutput
		var student models.Student
		var rating models.SessionRating

		err := rows.Scan(
			&result.SessionID, &result.Present, &result.Notes, &result.CreatedAt, &result.UpdatedAt,
			&student.ID, &student.FirstName, &student.LastName, &student.DOB, &student.TherapistID,
			&student.Grade, &student.IEP, &student.CreatedAt, &student.UpdatedAt,
			&rating.Level, &rating.Category, &rating.Description,
		)
		if err != nil {
			return nil, err
		}

		if existing, exists := sessionStudentsMap[student.ID.String()]; exists {
			// Student already exists
			if rating.Level != nil && rating.Category != nil {
				existing.Ratings = append(existing.Ratings, rating)
			}
		} else {
			result.Student = student
			result.Ratings = []models.SessionRating{}

			// Only add the rating if it's valid
			if rating.Level != nil && rating.Category != nil {
				result.Ratings = append(result.Ratings, rating)
			}

			sessionStudentsMap[student.ID.String()] = &result
		}
	}

	var sessionStudents []models.SessionStudentsOutput
	for _, ss := range sessionStudentsMap {
		sessionStudents = append(sessionStudents, *ss)
	}

	return sessionStudents, nil
}
