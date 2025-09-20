package schema

import (
	"context"
	"specialstandard/internal/models"

	"github.com/google/uuid"
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

func (r *ThemeRepository) GetThemes(ctx context.Context) ([]models.Theme, error) {
	query := `
        SELECT id, theme_name, month, year, created_at, updated_at
        FROM theme
        ORDER BY year DESC, month DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var themes []models.Theme
	for rows.Next() {
		var theme models.Theme
		if err := rows.Scan(
			&theme.ID,
			&theme.Name,
			&theme.Month,
			&theme.Year,
			&theme.CreatedAt,
			&theme.UpdatedAt,
		); err != nil {
			return nil, err
		}
		themes = append(themes, theme)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return themes, nil
}

func (r *ThemeRepository) GetThemeByID(ctx context.Context, id uuid.UUID) (*models.Theme, error) {
	theme := &models.Theme{}

	query := `
        SELECT id, theme_name, month, year, created_at, updated_at
        FROM theme
        WHERE id = $1`

	row := r.db.QueryRow(ctx, query, id)

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

func (r *ThemeRepository) UpdateTheme(ctx context.Context, id uuid.UUID, input *models.UpdateThemeInput) (*models.Theme, error) {
	// First get the existing theme
	existingTheme, err := r.GetThemeByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Apply updates
	if input.Name != nil {
		existingTheme.Name = *input.Name
	}
	if input.Month != nil {
		existingTheme.Month = *input.Month
	}
	if input.Year != nil {
		existingTheme.Year = *input.Year
	}

	query := `
        UPDATE theme
        SET theme_name = $1, month = $2, year = $3, updated_at = now()
        WHERE id = $4
        RETURNING id, theme_name, month, year, created_at, updated_at`

	row := r.db.QueryRow(ctx, query, existingTheme.Name, existingTheme.Month, existingTheme.Year, id)

	theme := &models.Theme{}
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

func (r *ThemeRepository) DeleteTheme(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM theme WHERE id = $1`
	
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func NewThemeRepository(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{
		db: db,
	}
}
