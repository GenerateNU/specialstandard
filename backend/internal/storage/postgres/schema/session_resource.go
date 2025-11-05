package schema

import (
	"context"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionResourceRepository struct {
	db *pgxpool.Pool
}

func (sr *SessionResourceRepository) PostSessionResource(ctx context.Context, sessionResource models.CreateSessionResource) (*models.SessionResource, error) {
	var newSessionResource models.SessionResource
	query := `INSERT INTO session_resource (session_id, resource_id)
				VALUES ($1, $2)
				RETURNING session_id, resource_id, created_at, updated_at`

	err := sr.db.QueryRow(ctx, query, sessionResource.SessionID, sessionResource.ResourceID).Scan(&newSessionResource.SessionID, &newSessionResource.ResourceID, &newSessionResource.CreatedAt, &newSessionResource.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &newSessionResource, nil
}

func (sr *SessionResourceRepository) DeleteSessionResource(ctx context.Context, sessionResource models.DeleteSessionResource) error {
	query := `DELETE FROM session_resource WHERE session_id = $1 AND resource_id = $2`
	_, err := sr.db.Exec(ctx, query, sessionResource.SessionID, sessionResource.ResourceID)
	if err != nil {
		return err
	}
	return nil
}

func (sr *SessionResourceRepository) GetResourcesBySessionID(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.Resource, error) {
	resources := make([]models.Resource, 0)
	query := `SELECT r.id, r.theme_id, r.grade_level, r.date, r.type, r.title, r.category, r.content, r.created_at, r.updated_at
				FROM session_resource sr
				JOIN resource r ON sr.resource_id = r.id
				WHERE sr.session_id = $1
				ORDER BY r.created_at ASC
				LIMIT $2 OFFSET $3`

	rows, err := sr.db.Query(ctx, query, sessionID, pagination.Limit, pagination.GetOffset())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Resource
		if err := rows.Scan(&r.ID, &r.ThemeID, &r.GradeLevel, &r.Date, &r.Type, &r.Title, &r.Category, &r.Content, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, err
		}
		resources = append(resources, r)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resources, nil
}

func NewSessionResourceRepository(db *pgxpool.Pool) *SessionResourceRepository {
	return &SessionResourceRepository{
		db,
	}
}
