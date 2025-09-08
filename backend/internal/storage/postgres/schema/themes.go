package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ThemeRepository struct {
	db *pgxpool.Pool
}

func (r *ThemeRepository) CreateTheme(ctx context.Context, input *models.CreateThemeInput) (*models.Theme, error) {
	// Create a Theme object to return
	theme := &models.Theme{}

	query := `
        INSERT INTO theme (theme_name, month, year)
        VALUES ($1, $2, $3)
        RETURNING id, theme_name, month, year, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, input.Name, input.Month, input.Year)

	// Scan into the theme object
	if err := row.Scan(
		&theme.ID,
		&theme.Name,
		&theme.Month,
		&theme.Year,
		&theme.CreatedAt,
		&theme.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return theme, nil
}

func NewThemeRepository(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{
		db: db,
	}
}
