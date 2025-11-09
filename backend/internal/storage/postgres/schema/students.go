package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

func (r *StudentRepository) GetStudents(ctx context.Context, grade *int, therapistID uuid.UUID, schoolID *int, name string, pagination utils.Pagination) ([]models.Student, error) {
	queryString := `
	SELECT s.id, s.first_name, s.last_name, s.dob, s.therapist_id, s.school_id, sch.name AS school_name, sch.district_id, s.grade, s.iep, s.created_at, s.updated_at
	FROM student s
	JOIN school sch ON s.school_id = sch.id
	WHERE 1 = 1`

	args := []interface{}{}
	argNum := 1

	// Add filters if provided (nil means no grade filter)
	if grade != nil {
		queryString += fmt.Sprintf(" AND s.grade = $%d", argNum)
		args = append(args, *grade)
		argNum++
	} else {
		queryString += " AND s.grade != -1" // Exclude graduated students by default
	}

	if therapistID != uuid.Nil {
		queryString += fmt.Sprintf(" AND s.therapist_id = $%d", argNum)
		args = append(args, therapistID)
		argNum++
	}

	if name != "" {
		queryString += fmt.Sprintf(" AND (s.first_name ILIKE $%d OR s.last_name ILIKE $%d OR CONCAT(s.first_name, ' ', s.last_name) ILIKE $%d)", argNum, argNum, argNum)
		searchPattern := "%" + name + "%"
		args = append(args, searchPattern)
		argNum++
	}

	if schoolID != nil {
		queryString += fmt.Sprintf(" AND s.school_id = $%d", argNum)
		args = append(args, *schoolID)
		argNum++
	}

	queryString += " ORDER BY sch.name ASC, s.first_name ASC, s.last_name ASC, s.dob ASC"

	// Add pagination
	queryString += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, pagination.Limit, pagination.GetOffset())

	rows, err := r.db.Query(ctx, queryString, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Student])
	if err != nil {
		return nil, err
	}
	return students, nil
}

func (r *StudentRepository) DeleteStudent(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM student
	WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *StudentRepository) UpdateStudent(ctx context.Context, student models.Student) (models.Student, error) {
	query := `
	UPDATE student
	SET first_name = $1, last_name = $2, dob = $3, therapist_id = $4, school_id = $5, grade = $6, iep = $7
	WHERE id = $8
	RETURNING id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at`

	var updatedStudent models.Student
	err := r.db.QueryRow(ctx, query,
		student.FirstName,
		student.LastName,
		student.DOB,
		student.TherapistID,
		student.SchoolID,
		student.Grade,
		student.IEP,
		student.ID,
	).Scan(
		&updatedStudent.ID,
		&updatedStudent.FirstName,
		&updatedStudent.LastName,
		&updatedStudent.DOB,
		&updatedStudent.TherapistID,
		&updatedStudent.SchoolID,
		&updatedStudent.Grade,
		&updatedStudent.IEP,
		&updatedStudent.CreatedAt,
		&updatedStudent.UpdatedAt,
	)
	return updatedStudent, err
}

func (r *StudentRepository) AddStudent(ctx context.Context, student models.Student) (models.Student, error) {
	query := `
	INSERT INTO student (first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	RETURNING id, first_name, last_name, dob, therapist_id, school_id, grade, iep, created_at, updated_at`

	var createdStudent models.Student
	err := r.db.QueryRow(ctx, query,
		student.FirstName,
		student.LastName,
		student.DOB,
		student.TherapistID,
		student.SchoolID,
		student.Grade,
		student.IEP,
	).Scan(
		&createdStudent.ID,
		&createdStudent.FirstName,
		&createdStudent.LastName,
		&createdStudent.DOB,
		&createdStudent.TherapistID,
		&createdStudent.SchoolID,
		&createdStudent.Grade,
		&createdStudent.IEP,
		&createdStudent.CreatedAt,
		&createdStudent.UpdatedAt,
	)

	return createdStudent, err
}

func (r *StudentRepository) GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error) {
	query := `
	SELECT s.id, s.first_name, s.last_name, s.dob, s.therapist_id, s.school_id, sch.name AS school_name, sch.district_id, s.grade, s.iep, s.created_at, s.updated_at
	FROM student s JOIN school sch ON s.school_id = sch.id
	WHERE s.id = $1`

	var student models.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.DOB,
		&student.TherapistID,
		&student.SchoolID,
		&student.SchoolName,
		&student.DistrictID,
		&student.Grade,
		&student.IEP,
		&student.CreatedAt,
		&student.UpdatedAt,
	)
	if err != nil {
		return models.Student{}, err
	}
	return student, nil
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{
		db,
	}
}

// This is our function to get all the sessions associated with a specific student from PostGres DB
func (r *StudentRepository) GetStudentSessions(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRepositoryRequest) ([]models.StudentSessionsOutput, error) {
	query := `
	SELECT ss.student_id, ss.present, ss.notes, ss.created_at, ss.updated_at,
	       s.id, s.start_datetime, s.end_datetime, s.therapist_id, s.notes, 
	       s.created_at, s.updated_at
	FROM session_student ss
	JOIN session s ON ss.session_id = s.id
	WHERE ss.student_id = $1`

	conditions := []string{}
	args := []interface{}{studentID}
	argCount := 2

	if filter != nil {
		// Date filtering - similar to sessions
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

		// Start and end date filtering
		if filter.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("s.start_datetime >= $%d", argCount))
			args = append(args, *filter.StartDate)
			argCount++
		}

		if filter.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("s.end_datetime <= $%d", argCount))
			args = append(args, *filter.EndDate)
			argCount++
		}

		// Attendance filtering
		if filter.Present != nil {
			conditions = append(conditions, fmt.Sprintf("ss.present = $%d", argCount))
			args = append(args, *filter.Present)
			argCount++
		}
	}

	// Add conditions to query
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY s.start_datetime ASC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pagination.Limit, pagination.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var studentSessions []models.StudentSessionsOutput
	for rows.Next() {
		var result models.StudentSessionsOutput
		var session models.Session

		err := rows.Scan(
			&result.StudentID, &result.Present, &result.Notes, &result.CreatedAt, &result.UpdatedAt,
			&session.ID, &session.StartDateTime, &session.EndDateTime, &session.TherapistID, &session.Notes,
			&session.CreatedAt, &session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		result.Session = session
		studentSessions = append(studentSessions, result)
	}

	return studentSessions, nil
}

func (r *StudentRepository) GetStudentRatings(ctx context.Context, studentID uuid.UUID, pagination utils.Pagination, filter *models.GetStudentSessionsRatingsRequest) ([]models.StudentSessionsWithRatingsOutput, error) {
	query := `
	SELECT ss.session_id, ss.student_id, s.start_datetime, sr.level, sr.category, sr.description
	FROM session_student ss
	JOIN session s ON ss.session_id = s.id
	LEFT JOIN session_rating sr ON ss.id = sr.session_student_id
	WHERE ss.student_id = $1`

	conditions := []string{}
	args := []interface{}{studentID}
	argCount := 2

	if filter != nil {
		// Date filtering
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

		// Start and end date filtering
		if filter.StartDate != nil {
			conditions = append(conditions, fmt.Sprintf("s.start_datetime >= $%d", argCount))
			args = append(args, *filter.StartDate)
			argCount++
		}

		if filter.EndDate != nil {
			conditions = append(conditions, fmt.Sprintf("s.end_datetime <= $%d", argCount))
			args = append(args, *filter.EndDate)
			argCount++
		}

		if filter.Present != nil {
			conditions = append(conditions, fmt.Sprintf("ss.present = $%d", argCount))
			args = append(args, *filter.Present)
			argCount++
		}
		if filter.Category != nil && len(*filter.Category) > 0 {
			conditions = append(conditions, fmt.Sprintf("sr.category = ANY($%d)", argCount))
			args = append(args, *filter.Category)
			argCount++
		}

	}

	// Add conditions to query
	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY s.start_datetime ASC"

	// Add pagination
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, pagination.Limit, pagination.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionMap := make(map[uuid.UUID]*models.StudentSessionsWithRatingsOutput)
	for rows.Next() {
		var sessionID, studentID uuid.UUID
		var sessionDate time.Time
		var category, level, description *string

		err := rows.Scan(
			&sessionID, &studentID, &sessionDate, &level, &category, &description,
		)
		if err != nil {
			return nil, err
		}

		if _, exists := sessionMap[sessionID]; !exists {
			sessionMap[sessionID] = &models.StudentSessionsWithRatingsOutput{
				SessionID:   sessionID,
				StudentID:   studentID,
				SessionDate: sessionDate,
				Ratings:     []models.SessionRating{},
			}
		}

		// Only append rating if it's not null
		if level != nil && category != nil {
			sessionMap[sessionID].Ratings = append(sessionMap[sessionID].Ratings, models.SessionRating{
				Level:       level,
				Category:    category,
				Description: description,
			})
		}
	}

	result := []models.StudentSessionsWithRatingsOutput{}
	for _, v := range sessionMap {
		result = append(result, *v)
	}

	return result, nil
}

func (r *StudentRepository) PromoteStudents(ctx context.Context, input models.PromoteStudentsInput) error {
	baseQuery := `UPDATE student
				  SET grade = CASE
					WHEN grade BETWEEN 0 AND 11 THEN (grade + 1)
					WHEN grade = 12 THEN -1
					ELSE grade
				  END
				  WHERE therapist_id = $1
						AND grade != -1`

	args := []interface{}{input.TherapistID}
	if len(input.ExcludedStudentIDs) > 0 {
		baseQuery += ` AND id != ALL($2)`
		args = append(args, input.ExcludedStudentIDs)
	}

	_, err := r.db.Exec(ctx, baseQuery, args...)
	if err != nil {
		return err
	}

	return nil
}
