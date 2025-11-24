package school

import (
	"specialstandard/internal/storage"
)

type Handler struct {
	schoolRepository storage.SchoolRepository
}

func NewHandler(schoolRepository storage.SchoolRepository) *Handler {
	return &Handler{
		schoolRepository: schoolRepository,
	}
}
