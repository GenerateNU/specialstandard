package session_resource

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	sessionResourceRepository storage.SessionResourceRepository
	validator                 *xvalidator.XValidator
}

func NewHandler(sessionResourceRepository storage.SessionResourceRepository) *Handler {
	return &Handler{
		sessionResourceRepository: sessionResourceRepository,
		validator:                 xvalidator.Validator,
	}
}
