package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStudentRepository struct {
	db *pgxpool.Pool
}

func NewSessionStudentRepository(db *pgxpool.Pool) *SessionStudentRepository {
	return &SessionStudentRepository{db: db}
}

func (r *SessionStudentRepository) CreateSessionStudent(ctx context.Context, input *models.CreateSessionStudentInput) (*models.SessionStudent, error) {
	query := `
	INSERT INTO session_student (session_id, student_id, present, notes)
	VALUES ($1, $2, $3, $4)
	RETURNING session_id, student_id, present, notes, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.SessionID, input.StudentID, input.Present, input.Notes)
	sessionStudent := &models.SessionStudent{}
	err := row.Scan(&sessionStudent.SessionID, &sessionStudent.StudentID, &sessionStudent.Present, &sessionStudent.Notes, &sessionStudent.CreatedAt, &sessionStudent.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return sessionStudent, nil
}
