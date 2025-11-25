package newsletter

import (
	"specialstandard/internal/s3_client"
	"specialstandard/internal/storage"
)

type Handler struct {
	Repo     storage.NewsletterRepository
	s3Client *s3_client.Client
}

func NewHandler(repo storage.NewsletterRepository, s3Client *s3_client.Client) *Handler {
	return &Handler{Repo: repo, s3Client: s3Client}
}
