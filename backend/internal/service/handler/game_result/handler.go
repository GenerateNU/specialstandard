package game_result

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	gameResultRepository storage.GameResultRepository
	validator            *xvalidator.XValidator
}

func NewHandler(gameResultRepository storage.GameResultRepository) *Handler {
	return &Handler{
		gameResultRepository: gameResultRepository,
		validator:            xvalidator.Validator,
	}
}
