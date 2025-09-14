package schema

import (
	"context"
	"fmt"
	"net/mail"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"strings"
	"time"

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
		return nil, err
	}

	return &therapist, nil
}

func (r *TherapistRepository) GetTherapists(ctx context.Context) ([]models.Therapist, error) {
	query := `
	SELECT id, first_name, last_name, email, active, created_at, updated_at
	FROM therapist`

	rows, err := r.db.Query(ctx, query)

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

	// AYEEE EMAIL VALIDATION !!!
	_, err := mail.ParseAddress(input.Email)

	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO therapist (first_name, last_name, email)
        VALUES ($1, $2, $3)
        RETURNING id, first_name, last_name, email, active, created_at, updated_At`

	row := r.db.QueryRow(ctx, query, input.First_name, input.Last_name, input.Email)

	// Scan into the therapist object
	if err := row.Scan(
		&therapist.ID,
		&therapist.First_name,
		&therapist.Last_name,
		&therapist.Email,
		&therapist.Active,
		&therapist.Created_at,
		&therapist.Updated_at,
	); err != nil {
		return nil, err
	}

	return therapist, nil
}

func (r *TherapistRepository) DeleteTherapist(ctx context.Context, therapistID string) (string, error) {
	query := `
	DELETE 
	FROM therapist
	WHERE id=$1`

	_, err := r.db.Exec(ctx, query, therapistID)

	if err != nil {
		return "", errs.BadRequest("Error querying database for given ID")
	}

	return "User Deleted Successfully", nil
}

// Here, we are just iterating through all of the potential changes, and updating the DB accordingly!
func (r *TherapistRepository) PatchTherapist(ctx context.Context, therapistID string, updatedValue *models.UpdateTherapist) (*models.Therapist, error) {
	query := `UPDATE therapist uc SET`
	updates := []string{}
	args := []interface{}{}
	argCount := 1

	if updatedValue.First_name != nil {
		updates = append(updates, fmt.Sprintf("first_name = $%d", argCount))
		args = append(args, *updatedValue.First_name)
		argCount++
	}

	if updatedValue.Last_name != nil {
		updates = append(updates, fmt.Sprintf("last_name = $%d", argCount))
		args = append(args, *updatedValue.Last_name)
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
		return nil, errs.NotFound("No fields given to update.")
	}

	updates = append(updates, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())
	argCount++

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
		return nil, errs.BadRequest("error querying database for given user ID")
	}

	return &therapist, nil
}

func NewTherapistRepository(db *pgxpool.Pool) *TherapistRepository {
	return &TherapistRepository{
		db,
	}
}
