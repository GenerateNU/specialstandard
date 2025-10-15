package resource

import (
	"specialstandard/internal/s3_client"
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	resourceRepository storage.ResourceRepository
	validator          *xvalidator.XValidator
	s3Client           *s3_client.Client
}

func NewHandler(resourceRepository storage.ResourceRepository, s3Client *s3_client.Client) *Handler {
	return &Handler{
		resourceRepository: resourceRepository,
		validator:          xvalidator.Validator,
		s3Client:           s3Client,
	}
}
