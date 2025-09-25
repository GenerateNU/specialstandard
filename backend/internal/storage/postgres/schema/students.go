package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

func (r *StudentRepository) GetStudents(ctx context.Context) ([]models.Student, error) {
	query := `
	SELECT id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at
	FROM student`

	rows, err := r.db.Query(ctx, query)
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
	SET first_name = $1, last_name = $2, dob = $3, therapist_id = $4, grade = $5, iep = $6
	WHERE id = $7
	RETURNING id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at`

	var updatedStudent models.Student
	err := r.db.QueryRow(ctx, query,
		student.FirstName,
		student.LastName,
		student.DOB,
		student.TherapistID,
		student.Grade,
		student.IEP,
		student.ID,
	).Scan(
		&updatedStudent.ID,
		&updatedStudent.FirstName,
		&updatedStudent.LastName,
		&updatedStudent.DOB,
		&updatedStudent.TherapistID,
		&updatedStudent.Grade,
		&updatedStudent.IEP,
		&updatedStudent.CreatedAt,
		&updatedStudent.UpdatedAt,
	)
	return updatedStudent, err
}

func (r *StudentRepository) AddStudent(ctx context.Context, student models.Student) (models.Student, error) {
	query := `
	INSERT INTO student (first_name, last_name, dob, therapist_id, grade, iep, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, NOW())
	RETURNING id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at`
	
	var createdStudent models.Student
	err := r.db.QueryRow(ctx, query,
		student.FirstName,
		student.LastName,
		student.DOB,         
		student.TherapistID,
		student.Grade,
		student.IEP,
	).Scan(
		&createdStudent.ID,
		&createdStudent.FirstName,
		&createdStudent.LastName,
		&createdStudent.DOB,
		&createdStudent.TherapistID,
		&createdStudent.Grade,
		&createdStudent.IEP,
		&createdStudent.CreatedAt,
		&createdStudent.UpdatedAt,
	)
	
	return createdStudent, err
}


func (r *StudentRepository) GetStudent(ctx context.Context, id uuid.UUID) (models.Student, error) {
	query := `
	SELECT id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at
	FROM student
	WHERE id = $1`
	
	var student models.Student
	err := r.db.QueryRow(ctx, query, id).Scan(
		&student.ID,
		&student.FirstName,
		&student.LastName,
		&student.DOB,
		&student.TherapistID,
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
func (r *StudentRepository) GetStudentSessions(ctx context.Context, studentID uuid.UUID) ([]models.StudentSessionsOutput, error) {
	query := `
	SELECT ss.student_id, ss.present, ss.notes, ss.created_at, ss.updated_at,
	       s.id, s.start_datetime, s.end_datetime, s.therapist_id, s.notes, 
	       s.created_at, s.updated_at
	FROM session_student ss
	JOIN session s ON ss.session_id = s.id
	WHERE ss.student_id = $1`

	rows, err := r.db.Query(ctx, query, studentID)
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