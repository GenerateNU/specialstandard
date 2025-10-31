package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/storage/dbinterface"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionStudentRepository struct {
	db *pgxpool.Pool
}

func (r *SessionStudentRepository) GetDB() *pgxpool.Pool {
	return r.db
}

func NewSessionStudentRepository(db *pgxpool.Pool) *SessionStudentRepository {
	return &SessionStudentRepository{db: db}
}

func (r *SessionStudentRepository) CreateSessionStudent(ctx context.Context, q dbinterface.Queryable, input *models.CreateSessionStudentInput) (*[]models.SessionStudent, error) {
	query := `INSERT INTO session_student (session_id, student_id, present, notes)
              VALUES `
	args := []interface{}{}
	argCount := 1

	for _, sessionID := range input.SessionIDs {
		for _, studentID := range input.StudentIDs {
			query += fmt.Sprintf(`($%d, $%d, $%d, $%d), `, argCount, argCount+1, argCount+2, argCount+3)
			args = append(args, sessionID, studentID, input.Present, input.Notes)
			argCount += 4
		}
	}

	// Removing the Trailing Comma + Space
	query = query[:len(query)-2]
	query += ` RETURNING session_id, student_id, present, notes, created_at, updated_at`

	rows, err := q.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessionStudents, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.SessionStudent])
	if err != nil {
		return nil, err
	}

	return &sessionStudents, nil
}

func (r *SessionStudentRepository) DeleteSessionStudent(ctx context.Context, input *models.DeleteSessionStudentInput) error {
	query := `DELETE FROM session_student WHERE session_id = $1 AND student_id = $2`
	_, err := r.db.Exec(ctx, query, input.SessionID, input.StudentID)
	return err
}

func (r *SessionStudentRepository) PatchSessionStudent(ctx context.Context, input *models.PatchSessionStudentInput) (*models.SessionStudent, error) {
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
