func (sr *SessionResourceRepository) GetResourcesBySessionID(ctx context.Context, sessionID uuid.UUID, pagination utils.Pagination) ([]models.Resource, error) {
    // enforce default pagination
    if pagination.Limit == 0 {
        pagination.Limit = 10
    }
    if pagination.Page == 0 {
        pagination.Page = 1
    }

    resources := make([]models.Resource, 0)

    query := `
        SELECT r.id, r.theme_id, r.grade_level, r.date, r.type, r.title, r.category, r.content, r.created_at, r.updated_at
        FROM session_resource sr
        JOIN resource r ON sr.resource_id = r.id
        WHERE sr.session_id = $1
        ORDER BY r.created_at DESC
        LIMIT $2 OFFSET $3
    `

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
