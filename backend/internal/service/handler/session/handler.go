package session

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	sessionRepository storage.SessionRepository
	validator         *xvalidator.XValidator
}

func NewHandler(sessionRepository storage.SessionRepository) *Handler {
	return &Handler{
		sessionRepository: sessionRepository,
		validator:         xvalidator.Validator,
	}
}
