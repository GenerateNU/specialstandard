package verification

import (
	"specialstandard/internal/storage"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/resend/resend-go/v3"
)

type Handler struct {
	verificationRepo storage.VerificationRepository
	db               *pgxpool.Pool
	resendClient     *resend.Client
	fromEmail        string
}

// Createing a new verification handler
func NewHandler(verificationRepo storage.VerificationRepository, db *pgxpool.Pool, resendApiKey, fromEmail string) *Handler {
	resendClient := resend.NewClient(resendApiKey)
	
	return &Handler{
		verificationRepo: verificationRepo,
		db:               db,
		resendClient:     resendClient,
		fromEmail:        fromEmail,
	}
}
