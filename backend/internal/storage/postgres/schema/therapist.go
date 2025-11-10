package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TherapistRepository struct {
	db *pgxpool.Pool
}

func (r *TherapistRepository) GetTherapistByID(ctx context.Context, therapistID string) (*models.Therapist, error) {
	query := `
	WITH therapist_data AS (
		SELECT 
			t.id, 
			t.first_name, 
			t.last_name,
			t.email, 
			t.active, 
			t.schools, 
			t.district_id, 
			d.name as district_name,
			t.created_at, 
			t.updated_at
		FROM therapist t
		LEFT JOIN district d ON t.district_id = d.id
		WHERE t.id = $1
	),
	school_data AS (
		SELECT 
			td.id as therapist_id,
			COALESCE(
				json_agg(
					json_build_object(
						'id', s.id,
						'name', s.name
					) ORDER BY s.name
				) FILTER (WHERE s.id IS NOT NULL), 
				'[]'::json
			) as schools_json
		FROM therapist_data td
		LEFT JOIN school s ON s.id = ANY(td.schools)
		GROUP BY td.id
	)
	SELECT 
		td.*,
		sd.schools_json
	FROM therapist_data td
	JOIN school_data sd ON td.id = sd.therapist_id`

	var therapist models.Therapist
	var schoolsJSON []byte
	var schools []int

	err := r.db.QueryRow(ctx, query, therapistID).Scan(
		&therapist.ID,
		&therapist.FirstName,
		&therapist.LastName,
		&therapist.Email,
		&therapist.Active,
		&schools,
		&therapist.DistrictID,
		&therapist.DistrictName,
		&therapist.CreatedAt,
		&therapist.UpdatedAt,
		&schoolsJSON,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errs.NotFound("Therapist not found with given ID")
		}
		return nil, err
	}

	therapist.Schools = schools

	// school names in a separate slice
	var schoolData []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	if err := json.Unmarshal(schoolsJSON, &schoolData); err != nil {
		return nil, err
	}

	schoolNames := make([]string, 0, len(schoolData))
	for _, s := range schoolData {
		schoolNames = append(schoolNames, s.Name)
	}
	therapist.SchoolNames = &schoolNames

	return &therapist, nil
}

func (r *TherapistRepository) GetTherapists(ctx context.Context, pagination utils.Pagination) ([]models.Therapist, error) {
	query := `
	SELECT t.id, t.first_name, t.last_name, t.email, t.active, t.schools, t.district_id, t.created_at, t.updated_at
	FROM therapist t
	ORDER BY first_name ASC, last_name ASC
	LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, pagination.Limit, pagination.GetOffset())

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Here i am using CollectExactlyOneRow because the DB should not have duplicate therapists!
	therapists, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Therapist])

	if err != nil {
		return nil, err
	}

	return therapists, nil
}

func (r *TherapistRepository) CreateTherapist(ctx context.Context, input *models.CreateTherapistInput) (*models.Therapist, error) {
	// Create a Therapist object to return
	therapist := &models.Therapist{}

	query := `
        INSERT INTO therapist (id, first_name, last_name, schools, district_id, email)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, first_name, last_name, schools, district_id, email, active, created_at, updated_At`

	row := r.db.QueryRow(ctx, query, input.ID, input.FirstName, input.LastName, input.Schools, input.DistrictID, input.Email)

	// Scan into the therapist object
	if err := row.Scan(
		&therapist.ID,
		&therapist.FirstName,
		&therapist.LastName,
		&therapist.Schools,
		&therapist.DistrictID,
		&therapist.Email,
		&therapist.Active,
		&therapist.CreatedAt,
		&therapist.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return therapist, nil
}

func (r *TherapistRepository) DeleteTherapist(ctx context.Context, therapistID string) error {
	query := `
	DELETE 
	FROM therapist
	WHERE id=$1`

	// db.Exec does not return an error if id isnt found, but
	// thats fine because our app will return 200 regardless if the ID exists or not!
	_, err := r.db.Exec(ctx, query, therapistID)

	// We will handle in the handler!
	if err != nil {
		return err
	}

	return nil
}

// Here, we are just iterating through all of the potential changes, and updating the DB accordingly!
func (r *TherapistRepository) PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error) {
	query := `UPDATE therapist uc SET`
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if updatedValue.FirstName != nil {
		updates = append(updates, fmt.Sprintf("first_name = $%d", argCount))
		args = append(args, *updatedValue.FirstName)
		argCount++
	}

	if updatedValue.LastName != nil {
		updates = append(updates, fmt.Sprintf("last_name = $%d", argCount))
		args = append(args, *updatedValue.LastName)
		argCount++
	}

	if updatedValue.Email != nil {
		updates = append(updates, fmt.Sprintf("email = $%d", argCount))
		args = append(args, *updatedValue.Email)
		argCount++
	}

	if updatedValue.Schools != nil {
		updates = append(updates, fmt.Sprintf("schools = $%d", argCount))
		args = append(args, *updatedValue.Schools)
		argCount++
	}

	if updatedValue.DistrictID != nil {
		updates = append(updates, fmt.Sprintf("district_id = $%d", argCount))
		args = append(args, *updatedValue.DistrictID)
		argCount++
	}

	if updatedValue.Active != nil {
		updates = append(updates, fmt.Sprintf("active = $%d", argCount))
		args = append(args, *updatedValue.Active)
		argCount++
	}

	if len(updates) == 0 {
		return nil, errs.BadRequest("No fields given to update.")
	}

	query += " " + strings.Join(updates, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, therapistID)

	query += " RETURNING id, first_name, last_name, email, schools, district_id, active, created_at, updated_At"

	rows, err := r.db.Query(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	therapist, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Therapist])

	if err != nil {
		return nil, errs.NotFound("error querying database for given user ID")
	}

	return &therapist, nil
}

func NewTherapistRepository(db *pgxpool.Pool) *TherapistRepository {
	return &TherapistRepository{
		db,
	}
}
