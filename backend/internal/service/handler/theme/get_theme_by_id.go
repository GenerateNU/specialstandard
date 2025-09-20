package theme

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetThemeByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
		slog.Error("Invalid UUID format for theme ID", "id", idStr, "error", err)
		return errs.BadRequest("Invalid UUID format")
	}

	theme, err := h.themeRepository.GetThemeByID(c.Context(), id)
	if err != nil {
		// Repository returns structured errors, just log and return them
		slog.Error("Database error getting theme", "id", id, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(theme)
}
