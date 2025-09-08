package theme

import "specialstandard/internal/storage"

type Handler struct {
	themeRepository storage.ThemeRepository
}

func NewHandler(themeRepository storage.ThemeRepository) *Handler {
	return &Handler{
		themeRepository,
	}
}
