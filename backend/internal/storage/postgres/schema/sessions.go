package schema

import (
	"context"
	"errors"
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

func (r *SessionRepository) GetSessions(ctx context.Context, pagination utils.Pagination, filter *models.GetSessionRepositoryRequest, therapistID uuid.UUID) ([]models.Session, error) {
	query := `
	SELECT s.id, s.session_name, s.start_datetime, s.end_datetime,
	       s.notes, s.location, s.created_at, s.updated_at,
	       s.session_parent_id,
		   sp.therapist_id,
	       sp.start_date, sp.end_date, sp.every_n_weeks, sp.days
	FROM session s
	INNER JOIN session_parent sp ON s.session_parent_id = sp.id
	`

	conditions := []string{}
	args := []interface{}{}
	argCount := 1

	// Therapist filter
	if therapistID != uuid.Nil {
		conditions = append(conditions, fmt.Sprintf("sp.therapist_id = $%d", argCount))
		args = append(args, therapistID)
		argCount++
	}

	if filter != nil {
		if filter.Month != nil && filter.Year != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(MONTH FROM s.start_datetime) = $%d AND EXTRACT(YEAR FROM s.start_datetime) = $%d", argCount, argCount+1))
			args = append(args, *filter.Month, *filter.Year)
			argCount += 2
		} else if filter.Year != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(YEAR FROM s.start_datetime) = $%d", argCount))
			args = append(args, *filter.Year)
			argCount++
		} else if filter.Month != nil {
			conditions = append(conditions, fmt.Sprintf("EXTRACT(MONTH FROM s.start_datetime) = $%d", argCount))
			args = append(args, *filter.Month)
			argCount++
		}

		if filter.StartTime != nil {
			conditions = append(conditions, fmt.Sprintf("s.start_datetime >= $%d", argCount))
			args = append(args, *filter.StartTime)
			argCount++
		}

		if filter.EndTime != nil {
			conditions = append(conditions, fmt.Sprintf("s.end_datetime <= $%d", argCount))
			args = append(args, *filter.EndTime)
			argCount++
		}

		if filter.StudentIDs != nil && len(*filter.StudentIDs) > 0 {
			for _, studentID := range *filter.StudentIDs {
				conditions = append(conditions, fmt.Sprintf(
					"EXISTS (SELECT 1 FROM session_student WHERE session_id = s.id AND student_id = $%d)",
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

	query += fmt.Sprintf(" ORDER BY s.start_datetime ASC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pagination.Limit, pagination.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var s models.Session
		var recurStart, recurEnd *time.Time
		var everyNWeeks *int
		var days []int

		if err := rows.Scan(
			&s.ID,
			&s.SessionName,
			&s.StartDateTime,
			&s.EndDateTime,
			&s.Notes,
			&s.Location,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.SessionParentID,
			&s.TherapistID,
			&recurStart,
			&recurEnd,
			&everyNWeeks,
			&days,
		); err != nil {
			return nil, err
		}

		// Populate repetition only if recur_start != recur_end
		if recurStart != nil && recurEnd != nil && !recurStart.Equal(*recurEnd) && everyNWeeks != nil {
			s.Repetition = &models.Repetition{
				RecurStart:  *recurStart,
				RecurEnd:    *recurEnd,
				EveryNWeeks: *everyNWeeks,
				Days:        days,
			}
		} else {
			s.Repetition = nil
		}

		sessions = append(sessions, s)
	}

	return sessions, nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (*models.Session, error) {
	query := `
	SELECT s.id, s.session_name, s.start_datetime, s.end_datetime,
	       s.notes, s.location, s.created_at, s.updated_at,
	       s.session_parent_id,
	       sp.start_date, sp.end_date, sp.every_n_weeks, sp.days
	FROM session s
	INNER JOIN session_parent sp ON s.session_parent_id = sp.id
	WHERE s.id = $1`

	row := r.db.QueryRow(ctx, query, id)

	var s models.Session
	var recurStart, recurEnd *time.Time
	var everyNWeeks *int
	var days []int

	if err := row.Scan(
		&s.ID,
		&s.SessionName,
		&s.StartDateTime,
		&s.EndDateTime,
		&s.Notes,
		&s.Location,
		&s.CreatedAt,
		&s.UpdatedAt,
		&s.SessionParentID,
		&recurStart,
		&recurEnd,
		&everyNWeeks,
		&days,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}

	if recurStart != nil && recurEnd != nil && !recurStart.Equal(*recurEnd) && everyNWeeks != nil {
		s.Repetition = &models.Repetition{
			RecurStart:  *recurStart,
			RecurEnd:    *recurEnd,
			EveryNWeeks: *everyNWeeks,
			Days:        days,
		}
	} else {
		s.Repetition = nil
	}

	return &s, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM session WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func addNWeeks(timestamp time.Time, nWeeks int) time.Time {
	return timestamp.Add(time.Duration(24*7*nWeeks) * time.Hour)
}

func (r *SessionRepository) PostSession(
	ctx context.Context,
	q dbinterface.Queryable,
	input *models.PostSessionInput,
) (*[]models.Session, error) {

	fmt.Printf("Posting session with data: %+v\n", input)
	fmt.Printf("Time details: %+v\n", input.StartTime)

	// Insert session_parent first
	var parentID uuid.UUID
	var parentStart, parentEnd time.Time

	if input.Repetition == nil {
		// Single session → make start=end so repetition=NULL in frontend
		parentStart = input.StartTime
		parentEnd = input.EndTime

		err := q.QueryRow(ctx,
			`INSERT INTO session_parent (start_date, end_date, every_n_weeks, days, therapist_id)
             VALUES ($1, $2, NULL, NULL, $3)
             RETURNING id`,
			parentStart, parentEnd, input.TherapistID,
		).Scan(&parentID)
		if err != nil {
			return nil, err
		}

	} else {
		rp := input.Repetition

		parentStart = rp.RecurStart
		parentEnd = rp.RecurEnd
		everyNWeeks := &rp.EveryNWeeks
		days := rp.Days

		err := q.QueryRow(ctx,
			`INSERT INTO session_parent (start_date, end_date, every_n_weeks, days, therapist_id)
             VALUES ($1, $2, $3, $4, $5)
             RETURNING id`,
			parentStart, parentEnd, everyNWeeks, days, input.TherapistID,
		).Scan(&parentID)
		if err != nil {
			return nil, err
		}
	}

	// Insert all session occurrences

	type insertedSession struct {
		ID    uuid.UUID
		Start time.Time
		End   time.Time
	}

	var sessionsInserted []insertedSession

	// Always insert the first session
	{

		fmt.Printf("Posting session with data: %+v\n", input)
		fmt.Printf("Time details: %+v\n", input.StartTime)

		// Add these right before the INSERT:
		fmt.Printf("About to insert START: %v, END: %v\n", input.StartTime, input.EndTime)
		fmt.Printf("START > END? %v\n", input.StartTime.After(input.EndTime))
		fmt.Printf("START == END? %v\n", input.StartTime.Equal(input.EndTime))

		var id uuid.UUID
		err := q.QueryRow(ctx,
			`INSERT INTO session (session_name, start_datetime, end_datetime, notes, location, session_parent_id)
             VALUES ($1, $2, $3, $4, $5, $6)
             RETURNING id, start_datetime, end_datetime`,
			input.SessionName, input.StartTime, input.EndTime,
			input.Notes, input.Location, parentID,
		).Scan(&id, &input.StartTime, &input.EndTime)
		if err != nil {
			return nil, err
		}

		sessionsInserted = append(sessionsInserted, insertedSession{
			ID:    id,
			Start: input.StartTime,
			End:   input.EndTime,
		})
	}

	// Generate repeating sessions only if repetition exists
	if input.Repetition != nil {
		rp := input.Repetition

		// We start from the week of the initial session start date.
		startOfInitialWeek := input.StartTime

		// For each N-week interval until rp.RecurEnd
		for wkStart := startOfInitialWeek; !wkStart.After(rp.RecurEnd); wkStart = addNWeeks(wkStart, rp.EveryNWeeks) {

			for _, dayIndex := range rp.Days {
				// Compute the occurrence for this weekday in the repeating week.
				occStart := setToWeekday(wkStart, dayIndex, input.StartTime)
				occEnd := setToWeekday(wkStart, dayIndex, input.EndTime)

				if occStart.Before(rp.RecurStart) {
					continue
				}

				// Skip if the occurrence is beyond the recurrence end
				if occStart.After(rp.RecurEnd) {
					continue
				}

				fmt.Printf("About to insert repeating session START: %v, END: %v\n", occStart, occEnd)
				fmt.Printf("START > END? %v\n", occStart.After(occEnd))
				fmt.Printf("START == END? %v\n", occStart.Equal(occEnd))

				var id uuid.UUID
				err := q.QueryRow(ctx,
					`INSERT INTO session (session_name, start_datetime, end_datetime, notes, location, session_parent_id)
                 VALUES ($1, $2, $3, $4, $5, $6)
                 RETURNING id, start_datetime, end_datetime`,
					input.SessionName, occStart, occEnd,
					input.Notes, input.Location, parentID,
				).Scan(&id, &occStart, &occEnd)
				if err != nil {
					return nil, err
				}

				sessionsInserted = append(sessionsInserted, insertedSession{
					ID:    id,
					Start: occStart,
					End:   occEnd,
				})
			}
		}
	}

	// Fetch full session objects using same logic as GetSessionByID (JOIN + scan repetition)

	sessions := make([]models.Session, 0, len(sessionsInserted))

	for _, sInserted := range sessionsInserted {

		row := q.QueryRow(ctx, `
            SELECT s.id, s.session_name, s.start_datetime, s.end_datetime,
                   s.notes, s.location, s.created_at, s.updated_at,
                   s.session_parent_id,
                   sp.start_date, sp.end_date, sp.every_n_weeks, sp.days
            FROM session s
            INNER JOIN session_parent sp ON s.session_parent_id = sp.id
            WHERE s.id = $1
        `, sInserted.ID)

		var s models.Session
		var recurStart, recurEnd *time.Time
		var enWeeks *int
		var d []int

		err := row.Scan(
			&s.ID,
			&s.SessionName,
			&s.StartDateTime,
			&s.EndDateTime,
			&s.Notes,
			&s.Location,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.SessionParentID,
			&recurStart,
			&recurEnd,
			&enWeeks,
			&d,
		)
		if err != nil {
			return nil, err
		}

		// Same repetition logic as GetSessionByID
		if recurStart != nil && recurEnd != nil && !recurStart.Equal(*recurEnd) && enWeeks != nil {
			s.Repetition = &models.Repetition{
				RecurStart:  *recurStart,
				RecurEnd:    *recurEnd,
				EveryNWeeks: *enWeeks,
				Days:        d,
			}
		} else {
			s.Repetition = nil
		}

		sessions = append(sessions, s)
	}

	return &sessions, nil
}

func setToWeekday(base time.Time, targetWeekday int, reference time.Time) time.Time {
	baseWeekday := int(base.Weekday()) // Sunday=0
	diff := targetWeekday - baseWeekday
	d := base.AddDate(0, 0, diff)

	// Preserve the hour/min/sec from the reference
	return time.Date(
		d.Year(), d.Month(), d.Day(),
		reference.Hour(), reference.Minute(), reference.Second(),
		reference.Nanosecond(), reference.Location(),
	)
}

func (r *SessionRepository) PatchSession(ctx context.Context, id uuid.UUID, input *models.PatchSessionInput) (*models.Session, error) {
	session := &models.Session{}

	query := `UPDATE session
				SET
					session_name = COALESCE($1, session_name),
			      start_datetime = COALESCE($2, start_datetime),
			      end_datetime = COALESCE($3, end_datetime),
					notes = COALESCE($4, notes),
					location = COALESCE($5, location)
				WHERE id = $6
				RETURNING id, session_name, start_datetime, end_datetime, notes, location, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.SessionName, input.StartTime, input.EndTime, input.Notes, input.Location, id)

	if err := row.Scan(
		&session.ID,
		&session.SessionName,
		&session.StartDateTime,
		&session.EndDateTime,
		&session.Notes,
		&session.Location,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	session.TherapistID = *input.TherapistID

	return session, nil
}

func (r *SessionRepository) DeleteRecurringSessions(ctx context.Context, id uuid.UUID) error {
	query := `
        DELETE FROM session
        WHERE session_parent_id = $1
        AND start_datetime > NOW();
    `

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *SessionRepository) GetSessionStudents(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination, therapistID uuid.UUID) ([]models.SessionStudentsOutput, error) {
	// Validate that therapistID is provided
	if therapistID == uuid.Nil {
		return nil, fmt.Errorf("therapist_id is required")
	}

	query := `
    SELECT ss.id, ss.session_id, ss.present, ss.notes, ss.created_at, ss.updated_at,
           s.id, s.first_name, s.last_name, s.dob, s.therapist_id, 
           s.grade, s.iep, s.created_at, s.updated_at, 
           sr.level, sr.category, sr.description
    FROM session_student ss
    JOIN student s ON ss.student_id = s.id
    LEFT JOIN session_rating sr ON ss.id = sr.session_student_id
    WHERE ss.session_id = $1
    AND s.therapist_id = $2
    AND s.grade != -1
    ORDER BY s.first_name, s.last_name
    LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(ctx, query, sessionID, therapistID, pagination.Limit, pagination.GetOffset())
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
			&result.SessionStudentID, &result.SessionID, &result.Present, &result.Notes, &result.CreatedAt, &result.UpdatedAt,
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

func NewSessionRepository(db *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		db,
	}
}
