package student

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	studentRepository storage.StudentRepository
	validator         *xvalidator.XValidator
}

func NewHandler(studentRepository storage.StudentRepository) *Handler {
	return &Handler{
		studentRepository: studentRepository,
		validator:         xvalidator.Validator,
	}
}
