package game_content

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	gameContentRepository storage.GameContentRepository
	validator             *xvalidator.XValidator
}

func NewHandler(gameContentRepository storage.GameContentRepository) *Handler {
	return &Handler{
		gameContentRepository: gameContentRepository,
		validator:             xvalidator.Validator,
	}
}
