package auth

import (
	"specialstandard/internal/config"
	"specialstandard/internal/storage"
)

type Handler struct {
	config                   config.Supabase
	therapistRepository      storage.TherapistRepository
	emailVerificationEnabled bool
}

type Credentials struct {
	Email      string  `json:"email"`
	Password   string  `json:"password"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	RememberMe bool    `json:"remember_me"`
}

func NewHandler(config config.Supabase, therapistRepository storage.TherapistRepository, emailVerificationEnabled bool) *Handler {
	return &Handler{
		config,
		therapistRepository,
		emailVerificationEnabled,
	}
}
