package schema

import (
	"context"
	"specialstandard/internal/models"

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
	query := `SELECT id, 
       				 category, 
       				 level, 
       		  		 (SELECT array_agg(opt) FROM 
       		  				 (SELECT opt FROM unnest(gc.options) AS opt ORDER BY random() LIMIT $3) AS sampled)
       		  		     	 AS options, 
       		  answer, 
       		  created_at, 
       		  updated_at 
			  FROM game_content gc
			  WHERE category = $1 AND level = $2;`

	rows, err := r.db.Query(ctx, query, req.Category, req.Level, req.Count-1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gameContents, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.GameContent])
	if err != nil {
		return nil, err
	}

	return gameContents, nil
}
