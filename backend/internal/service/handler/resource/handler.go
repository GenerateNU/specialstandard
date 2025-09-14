package resource

import "specialstandard/internal/storage"

type Handler struct {
	resourceRepository storage.ResourceRepository
}

func NewHandler(resourceRepository storage.ResourceRepository) *Handler {
	return &Handler{
		resourceRepository,
	}
}
