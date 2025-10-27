package session

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	sessionRepository        storage.SessionRepository
	sessionStudentRepository storage.SessionStudentRepository
	validator                *xvalidator.XValidator
}

func NewHandler(sessionRepository storage.SessionRepository, sessionStudentRepository storage.SessionStudentRepository) *Handler {
	return &Handler{
		sessionRepository:        sessionRepository,
		sessionStudentRepository: sessionStudentRepository,
		validator:                xvalidator.Validator,
	}
}
