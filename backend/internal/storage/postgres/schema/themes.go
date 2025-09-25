package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		return nil, errs.InternalServerError("Failed to create theme")
	}

	return theme, nil
}

func (r *ThemeRepository) GetThemes(ctx context.Context, pagination utils.Pagination) ([]models.Theme, error) {
	query := `
	SELECT id, theme_name, month, year, created_at, updated_at
	FROM theme
	ORDER BY year DESC, month DESC
	LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, pagination.Limit, pagination.GettOffset())

	if err != nil {
		return nil, errs.InternalServerError("Database connection error")
	}

	defer rows.Close()

	// Using CollectRows for simplicity
	themes, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Theme])

	if err != nil {
		return nil, errs.InternalServerError("Failed to retrieve themes")
	}

	return themes, nil
}

func (r *ThemeRepository) GetThemeByID(ctx context.Context, id uuid.UUID) (*models.Theme, error) {
	query := `
	SELECT id, theme_name, month, year, created_at, updated_at
	FROM theme
	WHERE id = $1`

	row, err := r.db.Query(ctx, query, id)

	if err != nil {
		return nil, errs.InternalServerError("Database connection error")
	}

	defer row.Close()

	// Using CollectExactlyOneRow
	theme, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.Theme])

	if err != nil {
		return nil, errs.NotFound("Error querying database for given ID")
	}

	return &theme, nil
}

func (r *ThemeRepository) PatchTheme(ctx context.Context, id uuid.UUID, input *models.UpdateThemeInput) (*models.Theme, error) {
	query := `UPDATE theme SET`
	updates := []string{}
	args := []interface{}{}
	argCount := 1

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

	if len(updates) == 0 {
		return nil, errs.BadRequest("No fields given to update.")
	}

	query += " " + strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	query += " RETURNING id, theme_name, month, year, created_at, updated_at"

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, errs.InternalServerError("Database connection error")
	}

	defer rows.Close()

	theme, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Theme])

	if err != nil {
		return nil, errs.NotFound("error querying database for given theme ID")
	}

	return &theme, nil
}

func (r *ThemeRepository) DeleteTheme(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE 
	FROM theme
	WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return errs.InternalServerError("Database connection error")
	}

	return nil
}

func NewThemeRepository(db *pgxpool.Pool) *ThemeRepository {
	return &ThemeRepository{
		db,
	}
}
