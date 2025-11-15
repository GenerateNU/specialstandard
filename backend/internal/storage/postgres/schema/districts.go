package schema

import (
	"context"
	"fmt"
	"specialstandard/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DistrictRepository struct {
	db *pgxpool.Pool
}

func NewDistrictRepository(db *pgxpool.Pool) *DistrictRepository {
	return &DistrictRepository{db: db}
}

// GetDistricts retrieves all districts
func (r *DistrictRepository) GetDistricts(ctx context.Context) ([]models.District, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM district
		ORDER BY name ASC
	`
	
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Error fetching districts from DB: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	
	var districts []models.District
	for rows.Next() {
		var district models.District
		err := rows.Scan(
			&district.ID,
			&district.Name,
			&district.CreatedAt,
			&district.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		districts = append(districts, district)
	}
	
	return districts, nil
}

// GetDistrictByID retrieves a single district by ID
func (r *DistrictRepository) GetDistrictByID(ctx context.Context, id int) (*models.District, error) {
	query := `
		SELECT id, name, created_at, updated_at 
		FROM district 
		WHERE id = $1
	`
	
	var district models.District
	err := r.db.QueryRow(ctx, query, id).Scan(
		&district.ID,
		&district.Name,
		&district.CreatedAt,
		&district.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &district, nil
}