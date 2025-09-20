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
		// Check if theme was not found (repository returns custom error string)
		if err.Error() == "theme not found" {
			slog.Error("Theme not found for deletion", "id", id, "error", err)
			return errs.NotFound("Theme not found")
		}
		// Other database errors
		slog.Error("Failed to delete theme", "id", id, "error", err)
		return errs.InternalServerError("Failed to delete theme")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Theme deleted successfully",
	})
}
