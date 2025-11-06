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
	query := `SELECT id, session_id, student_id, content_id, time_taken, completed, incorrect_tries, created_at, updated_at 
			  FROM game_result`

	var conditions []string
	var args []interface{}
	argCount := 1

	if inputQuery != nil {
		if inputQuery.SessionID != nil {
			conditions = append(conditions, fmt.Sprintf("session_id = $%d", argCount))
			args = append(args, inputQuery.SessionID)
			argCount++
		}

		if inputQuery.StudentID != nil {
			conditions = append(conditions, fmt.Sprintf("student_id = $%d", argCount))
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
	query := `INSERT INTO game_result (session_id, student_id, content_id, time_taken, completed, incorrect_tries)
			  VALUES ($1, $2, $3, $4, COALESCE($5, FALSE), COALESCE($6, 0))
			  RETURNING id, session_id, student_id, content_id, time_taken, completed, incorrect_tries, created_at, updated_at;`

	row := r.db.QueryRow(ctx, query, input.SessionID, input.StudentID, input.ContentID, input.TimeTaken, input.Completed, input.IncorrectTries)

	gameResult := &models.GameResult{}
	if err := row.Scan(
		&gameResult.ID,
		&gameResult.SessionID,
		&gameResult.StudentID,
		&gameResult.ContentID,
		&gameResult.TimeTaken,
		&gameResult.Completed,
		&gameResult.IncorrectTries,
		&gameResult.CreatedAt,
		&gameResult.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return gameResult, nil
}
