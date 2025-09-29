package schema

import (
	"context"
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
	SELECT id, first_name, last_name, email, active, created_at, updated_at
	FROM therapist
	WHERE id=$1`

	row, err := r.db.Query(ctx, query, therapistID)

	if err != nil {
		return nil, err
	}

	defer row.Close()

	// Here i am using CollectExactlyOneRow because the DB should not have duplicate therapists!
	therapist, err := pgx.CollectExactlyOneRow(row, pgx.RowToStructByName[models.Therapist])

	if err != nil {
		return nil, errs.NotFound("Error querying database for given ID")
	}

	return &therapist, nil
}

func (r *TherapistRepository) GetTherapists(ctx context.Context, pagination utils.Pagination) ([]models.Therapist, error) {
	query := `
	SELECT id, first_name, last_name, email, active, created_at, updated_at
	FROM therapist
	LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, pagination.Limit, pagination.GettOffset())

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
        INSERT INTO therapist (first_name, last_name, email)
        VALUES ($1, $2, $3)
        RETURNING id, first_name, last_name, email, active, created_at, updated_At`

	row := r.db.QueryRow(ctx, query, input.FirstName, input.LastName, input.Email)

	// Scan into the therapist object
	if err := row.Scan(
		&therapist.ID,
		&therapist.FirstName,
		&therapist.LastName,
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

	query += " RETURNING id, first_name, last_name, email, active, created_at, updated_At"

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
