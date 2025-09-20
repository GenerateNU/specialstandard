package theme

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) PatchTheme(c *fiber.Ctx) error {
	// Get ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)

	// Check if UUID is valid
	if err != nil {
		slog.Error("Invalid UUID format for theme update", "id", idStr, "error", err)
		return errs.BadRequest("Invalid UUID format")
	}

	var req models.UpdateThemeInput
	
	if err := c.BodyParser(&req); err != nil {
		slog.Error("Failed to parse theme update JSON", "error", err)
		return errs.InvalidJSON("Failed to parse theme data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(req); len(validationErrors) > 0 {
		slog.Error("Theme update validation failed", "errors", validationErrors)
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Update the theme
	updatedTheme, err := h.themeRepository.UpdateTheme(c.Context(), id, &req)
	if err != nil {
		// Check specific error types
		errStr := err.Error()
		switch {
		case errStr == "no fields provided to update":
			slog.Error("No fields provided for theme update", "id", id, "error", err)
			return errs.BadRequest("No fields provided to update")
		case strings.Contains(errStr, "no rows") || errStr == "sql: no rows in result set":
			slog.Error("Theme not found for update", "id", id, "error", err)
			return errs.NotFound("Theme not found")
		case strings.Contains(errStr, "foreign key"):
			slog.Error("Foreign key constraint error updating theme", "id", id, "error", err)
			return errs.BadRequest("Invalid reference to related data")
		case strings.Contains(errStr, "connection refused"):
			slog.Error("Database connection error updating theme", "id", id, "error", err)
			return errs.InternalServerError("Database connection error")
		default:
			slog.Error("Failed to update theme", "id", id, "error", err)
			return errs.InternalServerError("Failed to update theme")
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedTheme)
}
