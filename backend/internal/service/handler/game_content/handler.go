package game_content

import (
	"specialstandard/internal/s3_client"
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	gameContentRepository storage.GameContentRepository
	validator             *xvalidator.XValidator
	s3Client              *s3_client.Client
}

func NewHandler(gameContentRepository storage.GameContentRepository, s3Client *s3_client.Client) *Handler {
	return &Handler{
		gameContentRepository: gameContentRepository,
		validator:             xvalidator.Validator,
		s3Client:              s3Client,
	}
}
