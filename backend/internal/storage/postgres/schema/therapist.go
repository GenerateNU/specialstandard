package schema

import (
	"context"
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

func NewTherapistRepository(db *pgxpool.Pool) *TherapistRepository {
	return &TherapistRepository{
		db,
	}
}
