package theme

import (
	"specialstandard/internal/storage"
	"specialstandard/internal/xvalidator"
)

type Handler struct {
	themeRepository storage.ThemeRepository
	validator       *xvalidator.XValidator
}

func NewHandler(themeRepository storage.ThemeRepository) *Handler {
	return &Handler{
		themeRepository: themeRepository,
		validator:       xvalidator.Validator,
	}
}
