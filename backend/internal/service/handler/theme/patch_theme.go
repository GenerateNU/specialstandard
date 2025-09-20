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
	updatedTheme, err := h.themeRepository.PatchTheme(c.Context(), id, &req)
	if err != nil {
		// Repository now returns structured errors, so we can directly return them
		// or add additional context/logging as needed
		slog.Error("Failed to update theme", "id", id, "error", err)
		
		// Check for other database errors that might not be structured
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid reference to related data")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database connection error")
		default:
			// Return the error as-is if it's already a structured error from repository
			return err
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedTheme)
}
