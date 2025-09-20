package resource

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	resourceRepository storage.ResourceRepository
	validator          *xvalidator.XValidator
}

func NewHandler(resourceRepository storage.ResourceRepository) *Handler {
	return &Handler{
		resourceRepository: resourceRepository,
		validator:          xvalidator.Validator,
	}
}
