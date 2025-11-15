package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SchoolRepository struct {
	db *pgxpool.Pool
}

func NewSchoolRepository(db *pgxpool.Pool) *SchoolRepository {
	return &SchoolRepository{db: db}
}

// GetSchools retrieves all schools with district information
func (r *SchoolRepository) GetSchools(ctx context.Context) ([]models.School, error) {
	query := `
		SELECT 
			s.id, 
			s.name, 
			s.district_id, 
			s.created_at, 
			s.updated_at 
		FROM school s
		LEFT JOIN district d ON s.district_id = d.id
		ORDER BY s.name ASC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Error fetching schools from DB: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.School, error) {
		var s models.School
		err := row.Scan(
			&s.ID,
			&s.Name,
			&s.DistrictID,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		return s, err
	})
}

// GetSchoolsByDistrict retrieves all schools for a specific district
func (r *SchoolRepository) GetSchoolsByDistrict(ctx context.Context, districtID int) ([]models.School, error) {
	query := `
		SELECT id, name, district_id, created_at, updated_at 
		FROM school 
		WHERE district_id = $1
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(ctx, query, districtID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	return pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.School, error) {
		var s models.School
		err := row.Scan(
			&s.ID,
			&s.Name,
			&s.DistrictID,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		return s, err
	})
}