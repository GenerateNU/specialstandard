package schema

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
    db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
    return &AuthRepository{db: db}
}

func (r *AuthRepository) GetUserEmail(ctx context.Context, userID string) (string, error) {
    var email string
    query := `SELECT email FROM auth.users WHERE id = $1`
    err := r.db.QueryRow(ctx, query, userID).Scan(&email)
    return email, err
}

func (r *AuthRepository) MarkEmailVerified(ctx context.Context, userID string) error {
    query := `
        UPDATE auth.users 
        SET raw_user_meta_data = jsonb_set(
            COALESCE(raw_user_meta_data, '{}'::jsonb),
            '{email_verified}',
            'true'
        ),
        updated_at = NOW()
        WHERE id = $1
    `
    _, err := r.db.Exec(ctx, query, userID)
    return err
}