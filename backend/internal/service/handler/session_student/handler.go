package sessionstudent

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	sessionStudentRepository storage.SessionStudentRepository
	validator                *xvalidator.XValidator
}

func NewHandler(sessionStudentRepository storage.SessionStudentRepository) *Handler {
	return &Handler{
		sessionStudentRepository: sessionStudentRepository,
		validator:                xvalidator.Validator,
	}
}
