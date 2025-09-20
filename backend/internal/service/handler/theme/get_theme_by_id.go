package theme

import (
	"errors"
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
		// Check if it's a "no rows found" error using pgx's error constant
		if errors.Is(err, pgx.ErrNoRows) {
			slog.Error("Theme not found", "id", id, "error", err)
			return errs.NotFound("Theme not found")
		}
		// For all other database errors, return internal server error without exposing details
		slog.Error("Database error getting theme", "id", id, "error", err)
		return errs.InternalServerError("Database error")
	}

	return c.Status(fiber.StatusOK).JSON(theme)
}
