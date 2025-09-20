package theme

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTheme(c *fiber.Ctx) error {
	var theme models.CreateThemeInput

	if err := c.BodyParser(&theme); err != nil {
		slog.Error("Failed to parse theme JSON", "error", err)
		return errs.InvalidJSON("Failed to parse theme data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(theme); len(validationErrors) > 0 {
		slog.Error("Theme validation failed", "errors", validationErrors)
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	createdTheme, err := h.themeRepository.CreateTheme(c.Context(), &theme)
	if err != nil {
		// Specific error handling with custom messages
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			slog.Error("Foreign key constraint error creating theme", "error", err)
			return errs.BadRequest("Invalid reference to related data")
		case strings.Contains(errStr, "connection refused"):
			slog.Error("Database connection error creating theme", "error", err)
			return errs.InternalServerError("Database connection error")
		default:
			slog.Error("Failed to create theme", "error", err)
			return errs.InternalServerError("Failed to create theme")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(createdTheme)
}
