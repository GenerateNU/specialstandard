package therapist

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	therapistRepository storage.TherapistRepository
	validator       *xvalidator.XValidator
}

func NewHandler(therapistRepository storage.TherapistRepository) *Handler {
	return &Handler{
		therapistRepository: therapistRepository,
		validator:       xvalidator.Validator,
	}
}
