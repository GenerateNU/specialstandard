package schema

import (
	"context"
	"database/sql"
	"specialstandard/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type NewsletterRepository struct {
	db *pgxpool.Pool
}

func NewNewsletterRepository(db *pgxpool.Pool) *NewsletterRepository {
	return &NewsletterRepository{db: db}
}

func (r *NewsletterRepository) GetNewsletterByDate(ctx context.Context, date time.Time) (*models.Newsletter, error) {
	var n models.Newsletter
	query := `SELECT id, start_date, end_date, s3_url FROM newsletter WHERE start_date <= $1 AND end_date >= $1 LIMIT 1`
	row := r.db.QueryRow(ctx, query, date)
	err := row.Scan(&n.ID, &n.StartDate, &n.EndDate, &n.S3URL)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &n, nil
}
