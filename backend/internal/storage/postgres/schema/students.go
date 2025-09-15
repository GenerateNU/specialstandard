package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/google/uuid"
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

func (r *StudentRepository) UpdateStudent(ctx context.Context, student models.Student) error {
	query := `
	UPDATE student
	SET first_name = $1, last_name = $2, dob = $3, therapist_id = $4, grade = $5, iep = $6, updated_at = NOW()
	WHERE id = $7`

	_, err := r.db.Exec(ctx, query,
		student.FirstName,
		student.LastName,
		student.DOB,            // Changed from DateOfBirth
		student.TherapistID,
		student.Grade,
		student.IEP,
		student.ID,
	)
	return err
}

func (r *StudentRepository) AddStudent(ctx context.Context, student models.Student) error {
	query := `
	INSERT INTO student (id, first_name, last_name, dob, therapist_id, grade, iep, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())`
	
	_, err := r.db.Exec(ctx, query,
		student.ID,
		student.FirstName,
		student.LastName,
		student.DOB,
		student.TherapistID,
		student.Grade,
		student.IEP,
	)
	return err
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
		&student.DOB,           // Changed from DateOfBirth
		&student.TherapistID,
		&student.Grade,
		&student.IEP,
		&student.CreatedAt,     // Added
		&student.UpdatedAt,     // Added
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
