package theme

import (
	"log/slog"
	"specialstandard/internal/errs"
	"strings"

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
		// Theme not found
		if strings.Contains(err.Error(), "no rows") || err.Error() == "sql: no rows in result set" {
			slog.Error("Theme not found", "id", id, "error", err)
			return errs.NotFound("Theme not found")
		}
		// Some other error
		slog.Error("Database error getting theme", "id", id, "error", err)
		return errs.InternalServerError("Database error")
	}

	return c.Status(fiber.StatusOK).JSON(theme)
}
