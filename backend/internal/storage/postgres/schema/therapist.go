package schema

import (
	"context"
	"net/mail"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"

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

func NewTherapistRepository(db *pgxpool.Pool) *TherapistRepository {
	return &TherapistRepository{
		db,
	}
}
