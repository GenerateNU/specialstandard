package schema

import (
	"context"
	"fmt"
	"log/slog"
	"specialstandard/internal/models"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GameContentRepository struct {
	db *pgxpool.Pool
}

func NewGameContentRepository(db *pgxpool.Pool) *GameContentRepository {
	return &GameContentRepository{
		db,
	}
}

func (r *GameContentRepository) GetGameContents(ctx context.Context, req models.GetGameContentRequest) ([]models.GameContent, error) {
	query := `SELECT id, theme_id, week, category, question_type, difficulty_level, question, 
             (SELECT array_agg(opt) 
              	FROM (SELECT opt FROM unnest(gc.options) AS opt ORDER BY random() LIMIT $1) AS sampled)
              	AS options,
    		 answer, created_at, updated_at
       	     FROM game_content gc`

	var conditions []string
	var args []interface{}
	args = append(args, *req.WordsCount-1)
	argCount := 2

	if req.ThemeID != nil {
		conditions = append(conditions, fmt.Sprintf("theme_id = $%d", argCount))
		args = append(args, *req.ThemeID)
		argCount++
	}
	if req.Category != nil {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argCount))
		args = append(args, *req.Category)
		argCount++
	}
	if req.QuestionType != nil {
		conditions = append(conditions, fmt.Sprintf("question_type = $%d", argCount))
		args = append(args, *req.QuestionType)
		argCount++
	}
	if req.DifficultyLevel != nil {
		conditions = append(conditions, fmt.Sprintf("difficulty_level = $%d", argCount))
		args = append(args, *req.DifficultyLevel)
		argCount++
	}

	if len(conditions) > 0 {
		query += ` WHERE ` + strings.Join(conditions, " AND ")
	}
	query += fmt.Sprintf(` ORDER BY random() LIMIT $%d `, argCount)
	args = append(args, *req.QuestionCount)
	// argCount++
	// (Commented out because linter says "ineffectual" but if we add arguments below this, PLEASE UNCOMMENT THAT LINE! Otherwise, number will be off by one.

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		slog.Error("Failed to get game contents", "error", err)
		return nil, err
	}

	gameContents, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.GameContent])
	if err != nil {
		return nil, err
	}

	return gameContents, nil
}
