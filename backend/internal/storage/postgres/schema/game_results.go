package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameResultRepository struct {
	db *pgxpool.Pool
}

func NewGameResultRepository(db *pgxpool.Pool) *GameResultRepository {
	return &GameResultRepository{
		db,
	}
}

func (r *GameResultRepository) GetGameResults(ctx context.Context, inputQuery *models.GetGameResultQuery, pagination utils.Pagination) ([]models.GameResult, error) {
	query := `SELECT gr.id, gr.session_student_id, gr.content_id, gr.time_taken_sec, gr.completed,
       					gr.count_of_incorrect_attempts, gr.incorrect_attempts, gr.created_at, gr.updated_at
			  FROM game_result gr JOIN session_student ss ON gr.session_student_id = ss.id`

	var conditions []string
	var args []interface{}
	argCount := 1

	if inputQuery != nil {
		if inputQuery.SessionID != nil {
			conditions = append(conditions, fmt.Sprintf("ss.session_id = $%d", argCount))
			args = append(args, inputQuery.SessionID)
			argCount++
		}

		if inputQuery.StudentID != nil {
			conditions = append(conditions, fmt.Sprintf("ss.student_id = $%d", argCount))
			args = append(args, inputQuery.StudentID)
			argCount++
		}
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT $%d OFFSET $%d`, argCount, argCount+1)
	args = append(args, pagination.Limit, pagination.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	gameResults, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.GameResult])
	if err != nil {
		return nil, err
	}

	return gameResults, nil
}

func (r *GameResultRepository) PostGameResult(ctx context.Context, input models.PostGameResult) (*models.GameResult, error) {
	query := `INSERT INTO game_result (session_student_id, content_id, time_taken_sec, completed, count_of_incorrect_attempts, incorrect_attempts)
			  VALUES ($1, $2, $3, COALESCE($4, FALSE), COALESCE($5, 0), COALESCE($6, '{}'))
			  RETURNING id, student_session_id, content_id, time_taken_sec, completed, count_of_incorrect_attempts, incorrect_attempts, created_at, updated_at;`

	row := r.db.QueryRow(ctx, query, input.SessionStudentID, input.ContentID, input.TimeTakenSec,
		input.Completed, input.CountIncorrectAttempts, input.IncorrectAttempts)

	gameResult := &models.GameResult{}
	if err := row.Scan(
		&gameResult.ID,
		&gameResult.SessionStudentID,
		&gameResult.ContentID,
		&gameResult.TimeTakenSec,
		&gameResult.Completed,
		&gameResult.CountIncorrectAttempts,
		&gameResult.IncorrectAttempts,
		&gameResult.CreatedAt,
		&gameResult.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return gameResult, nil
}
