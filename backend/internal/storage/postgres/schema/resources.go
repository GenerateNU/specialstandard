package schema

import (
	"context"
	"fmt"
	"log/slog"
	"specialstandard/internal/utils"
	"strings"
	"time"

	"specialstandard/internal/errs"
	"specialstandard/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResourceRepository struct {
	db *pgxpool.Pool
}

func NewResourceRepository(db *pgxpool.Pool) *ResourceRepository {
	return &ResourceRepository{
		db: db,
	}
}

func (r *ResourceRepository) GetResources(ctx context.Context, themeID uuid.UUID, gradeLevel, resType, title, category, content, themeName string, date *time.Time, themeMonth, themeYear int, pagination utils.Pagination) ([]models.ResourceWithTheme, error) {
	var resources []models.ResourceWithTheme
	queryString := "SELECT r.id, r.theme_id, r.grade_level, r.date, r.type, r.title, r.category, r.content, r.created_at, r.updated_at, t.theme_name, t.month, t.year, t.created_at, t.updated_at FROM resource r JOIN theme t ON r.theme_id = t.id WHERE 1=1"
	args := []interface{}{}
	argNum := 1

	if themeID != uuid.Nil {
		queryString += fmt.Sprintf(" AND theme_id = $%d::uuid", argNum)
		args = append(args, themeID)
		argNum++
	}
	if gradeLevel != "" {
		queryString += fmt.Sprintf(" AND grade_level = $%d", argNum)
		args = append(args, gradeLevel)
		argNum++
	}
	if resType != "" {
		queryString += fmt.Sprintf(" AND type = $%d", argNum)
		args = append(args, resType)
		argNum++
	}
	if title != "" {
		queryString += fmt.Sprintf(" AND title = $%d", argNum)
		args = append(args, title)
		argNum++
	}
	if category != "" {
		queryString += fmt.Sprintf(" AND category = $%d", argNum)
		args = append(args, category)
		argNum++
	}
	if content != "" {
		queryString += fmt.Sprintf(" AND content = $%d", argNum)
		args = append(args, content)
		argNum++
	}
	if date != nil {
		queryString += fmt.Sprintf(" AND date = $%d", argNum)
		args = append(args, date)
		argNum++
	}
	if themeName != "" {
		queryString += fmt.Sprintf(" AND t.theme_name ILIKE $%d", argNum)
		args = append(args, "%"+themeName+"%")
		argNum++
	}
	if themeMonth != 0 {
		queryString += fmt.Sprintf(" AND t.month = $%d", argNum)
		args = append(args, themeMonth)
		argNum++
	}
	if themeYear != 0 {
		queryString += fmt.Sprintf(" AND t.year = $%d", argNum)
		args = append(args, themeYear)
		argNum++
	}

	queryString += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, pagination.Limit, pagination.GettOffset())

	rows, err := r.db.Query(ctx, queryString, args...)
	if err != nil {
		slog.Error("Error occurred", "error", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var resource models.ResourceWithTheme
		err := rows.Scan(
			&resource.ID,
			&resource.ThemeID,
			&resource.GradeLevel,
			&resource.Date,
			&resource.Type,
			&resource.Title,
			&resource.Category,
			&resource.Content,
			&resource.CreatedAt,
			&resource.UpdatedAt,
			&resource.Theme.Name,
			&resource.Theme.Month,
			&resource.Theme.Year,
			&resource.Theme.CreatedAt,
			&resource.Theme.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		resources = append(resources, resource)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return resources, nil
}

func (r *ResourceRepository) GetResourceByID(ctx context.Context, id uuid.UUID) (*models.ResourceWithTheme, error) {
	var resource models.ResourceWithTheme
	err := r.db.QueryRow(ctx, "SELECT r.id, r.theme_id, r.grade_level, r.date, r.type, r.title, r.category, r.content, r.created_at, r.updated_at, t.theme_name, t.month, t.year, t.created_at, t.updated_at FROM resource r JOIN theme t ON r.theme_id = t.id WHERE r.id = $1", id.String()).Scan(
		&resource.ID,
		&resource.ThemeID,
		&resource.GradeLevel,
		&resource.Date,
		&resource.Type,
		&resource.Title,
		&resource.Category,
		&resource.Content,
		&resource.CreatedAt,
		&resource.UpdatedAt,
		&resource.Theme.Name,
		&resource.Theme.Month,
		&resource.Theme.Year,
		&resource.Theme.CreatedAt,
		&resource.Theme.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &resource, nil
}

func (r *ResourceRepository) UpdateResource(ctx context.Context, id uuid.UUID, resourceBody models.UpdateResourceBody) (*models.Resource, error) {
	var updatedResource models.Resource
	setFields := []string{}
	args := []interface{}{}
	argNum := 1

	if resourceBody.ThemeID != nil {
		setFields = append(setFields, "theme_id = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.ThemeID)
		argNum++
	}
	if resourceBody.GradeLevel != nil {
		setFields = append(setFields, "grade_level = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.GradeLevel)
		argNum++
	}
	if resourceBody.Date != nil {
		setFields = append(setFields, "date = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.Date)
		argNum++
	}
	if resourceBody.Type != nil {
		setFields = append(setFields, "type = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.Type)
		argNum++
	}
	if resourceBody.Title != nil {
		setFields = append(setFields, "title = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.Title)
		argNum++
	}
	if resourceBody.Category != nil {
		setFields = append(setFields, "category = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.Category)
		argNum++
	}
	if resourceBody.Content != nil {
		setFields = append(setFields, "content = $"+fmt.Sprint(argNum))
		args = append(args, *resourceBody.Content)
		argNum++
	}

	if len(setFields) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setFields = append(setFields, "updated_at = $"+fmt.Sprint(argNum))
	args = append(args, time.Now())
	argNum++

	query := "UPDATE resource SET " + strings.Join(setFields, ", ") + " WHERE id = $" + fmt.Sprint(argNum) + " RETURNING id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at"
	args = append(args, id.String())

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&updatedResource.ID,
		&updatedResource.ThemeID,
		&updatedResource.GradeLevel,
		&updatedResource.Date,
		&updatedResource.Type,
		&updatedResource.Title,
		&updatedResource.Category,
		&updatedResource.Content,
		&updatedResource.CreatedAt,
		&updatedResource.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &updatedResource, nil
}

func (r *ResourceRepository) CreateResource(ctx context.Context, resourceBody models.ResourceBody) (*models.Resource, error) {
	var newResource models.Resource
	err := r.db.QueryRow(ctx, "INSERT INTO resource (theme_id, grade_level, date, type, title, category, content) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, theme_id, grade_level, date, type, title, category, content, created_at, updated_at",
		resourceBody.ThemeID,
		resourceBody.GradeLevel,
		resourceBody.Date,
		resourceBody.Type,
		resourceBody.Title,
		resourceBody.Category,
		resourceBody.Content,
	).Scan(&newResource.ID, &newResource.ThemeID, &newResource.GradeLevel, &newResource.Date, &newResource.Type, &newResource.Title, &newResource.Category, &newResource.Content, &newResource.CreatedAt, &newResource.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "23503") {
			return nil, errs.InvalidRequestData(map[string]string{"theme_id": "invalid theme"})
		}
		return nil, errs.InternalServerError(err.Error())
	}
	return &newResource, nil
}

func (r *ResourceRepository) DeleteResource(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM resource WHERE id = $1", id)
	if err != nil {
		return errs.InternalServerError(err.Error())
	}

	return nil
}
