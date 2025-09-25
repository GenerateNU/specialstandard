package sessionstudent

import "specialstandard/internal/storage"

type Handler struct {
	sessionStudentRepository storage.SessionStudentRepository
}

func NewHandler(sessionStudentRepository storage.SessionStudentRepository) *Handler {
	return &Handler{
		sessionStudentRepository,
	}
}
