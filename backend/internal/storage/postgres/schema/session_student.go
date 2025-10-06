package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStudentRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewSessionStudentRepository(db *pgxpool.Pool) *SessionStudentRepositoryImpl {
	return &SessionStudentRepositoryImpl{db: db}
}

func (r *SessionStudentRepositoryImpl) CreateSessionStudent(ctx context.Context, input *models.CreateSessionStudentInput) (*models.SessionStudent, error) {
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

func (r *SessionStudentRepositoryImpl) DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error {
	query := `DELETE FROM session_student WHERE session_id = $1 AND student_id = $2`
	_, err := r.db.Exec(ctx, query, input.SessionID, input.StudentID)
	return err
}

func (r *SessionStudentRepositoryImpl) PatchSessionStudent(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, error) {
	sessionStudent := &models.SessionStudent{}

	query := `UPDATE session_student
				SET
					present = COALESCE($1, present),
					notes = COALESCE($2, notes)
				WHERE session_id = $3 AND student_id = $4
				RETURNING session_id, student_id, present, notes, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.Present, input.Notes, input.SessionID, input.StudentID)

	if err := row.Scan(
		&sessionStudent.SessionID,
		&sessionStudent.StudentID,
		&sessionStudent.Present,
		&sessionStudent.Notes,
		&sessionStudent.CreatedAt,
		&sessionStudent.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return sessionStudent, nil
}
