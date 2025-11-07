package schema

import (
	"context"
	"specialstandard/internal/models"

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
	query := `SELECT id, category, level, 
       		  		 (SELECT array_agg(opt) FROM 
       		  				 (SELECT opt FROM unnest(gc.options) AS opt ORDER BY random() LIMIT $3) AS sampled)
       		  		     	 AS options, 
       		  answer, created_at, updated_at 
			  FROM game_content gc
			  WHERE category = $1 AND level = $2;`

	row := r.db.QueryRow(ctx, query, req.Category, req.DifficultyLevel, req.Count-1)

	var gc models.GameContent
	if err := row.Scan(
		&gc.ID, &gc.Category, &gc.DifficultyLevel, &gc.Options, &gc.Answer, &gc.CreatedAt, &gc.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &gc, nil
}
