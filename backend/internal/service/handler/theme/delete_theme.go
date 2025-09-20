package theme

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) DeleteTheme(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	
	// Check if UUID is valid
	if err != nil {
		slog.Error("Invalid UUID format for theme deletion", "id", idStr, "error", err)
		return errs.BadRequest("Invalid UUID format")
	}
	
	if err := h.themeRepository.DeleteTheme(c.Context(), id); err != nil {
		// Repository returns structured errors, just log and return them
		slog.Error("Failed to delete theme", "id", id, "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Theme deleted successfully",
	})
}
