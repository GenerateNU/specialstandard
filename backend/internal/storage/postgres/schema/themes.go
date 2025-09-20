package schema

import (
	"context"
	"errors"
	"fmt"
	"specialstandard/internal/models"
	"strings"

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
	// Build dynamic SQL query - only update provided fields
	query := `UPDATE theme SET`
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	// Only add fields that are provided
	if input.Name != nil {
		updates = append(updates, fmt.Sprintf("theme_name = $%d", argCount))
		args = append(args, *input.Name)
		argCount++
	}

	if input.Month != nil {
		updates = append(updates, fmt.Sprintf("month = $%d", argCount))
		args = append(args, *input.Month)
		argCount++
	}

	if input.Year != nil {
		updates = append(updates, fmt.Sprintf("year = $%d", argCount))
		args = append(args, *input.Year)
		argCount++
	}

	// Validate that at least one field was provided
	if len(updates) == 0 {
		return nil, errors.New("no fields provided to update")
	}

	// Complete the query
	query += " " + strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)
	query += " RETURNING id, theme_name, month, year, created_at, updated_at"

	// Execute single atomic query
	row := r.db.QueryRow(ctx, query, args...)

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
	
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	
	// Check if any rows were actually deleted
	if result.RowsAffected() == 0 {
		return errors.New("theme not found")
	}
	
	return nil
}

func NewThemeRepository(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{
		db: db,
	}
}
