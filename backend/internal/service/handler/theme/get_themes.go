package theme

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetThemes(c *fiber.Ctx) error {
	themes, err := h.themeRepository.GetThemes(c.Context())
	if err != nil {
		slog.Error("Failed to fetch themes", "error", err)
		return errs.InternalServerError("Failed to fetch themes")
	}

	return c.Status(fiber.StatusOK).JSON(themes)
}
