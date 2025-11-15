package schema

import (
	"context"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VerificationRepository struct {
	db *pgxpool.Pool
}

func NewVerificationRepository(db *pgxpool.Pool) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) CreateVerificationCode(ctx context.Context, code models.VerificationCode) error {
	// When creating, 'used' defaults to false
	query := `
		INSERT INTO verification_codes (user_id, code, expires_at, created_at, used)
		VALUES ($1, $2, $3, $4, false)
	`

	_, err := r.db.Exec(ctx, query, code.UserID, code.Code, code.ExpiresAt, code.CreatedAt)
	if err != nil {
		return errs.InternalServerError("Failed to create verification code")
	}

	return nil
}

func (r *VerificationRepository) VerifyCode(ctx context.Context, userID, code string) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	// Check if code exists and is valid - 'used' is boolean
	var verificationCode models.VerificationCode
	query := `
		SELECT id, user_id, code, expires_at, used
		FROM verification_codes
		WHERE user_id = $1 AND code = $2 AND used = false
		FOR UPDATE
	`

	err = tx.QueryRow(ctx, query, userID, code).Scan(
		&verificationCode.ID,
		&verificationCode.UserID,
		&verificationCode.Code,
		&verificationCode.ExpiresAt,
		&verificationCode.Used,  // Now it's a boolean
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	// Check if code has expired
	if time.Now().After(verificationCode.ExpiresAt) {
		return false, nil
	}

	// Mark code as used - set boolean to true
	updateQuery := `
		UPDATE verification_codes
		SET used = true
		WHERE id = $1
	`

	_, err = tx.Exec(ctx, updateQuery, verificationCode.ID)
	if err != nil {
		return false, err
	}

	if err = tx.Commit(ctx); err != nil {
		return false, err
	}

	return true, nil
}

func (r *VerificationRepository) CleanupExpiredCodes(ctx context.Context) error {
	// Delete codes that are expired or already used
	query := `
		DELETE FROM verification_codes
		WHERE expires_at < $1 OR used = true
	`

	_, err := r.db.Exec(ctx, query, time.Now())
	return err
}