package student

import "specialstandard/internal/storage"

type Handler struct {
	studentRepository storage.StudentRepository
}

func NewHandler(studentRepository storage.StudentRepository) *Handler {
	return &Handler{
		studentRepository,
	}
}