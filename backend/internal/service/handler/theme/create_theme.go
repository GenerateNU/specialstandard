package theme

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTheme(c *fiber.Ctx) error {
	var theme models.CreateThemeInput

	if err := c.BodyParser(&theme); err != nil {
		return errs.InvalidJSON("Failed to parse theme data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(theme); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	createdTheme, err := h.themeRepository.CreateTheme(c.Context(), &theme)
	if err != nil {
		// Specific error handling with custom messages
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid reference to related data")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database connection error")
		default:
			return errs.InternalServerError("Failed to create theme")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(createdTheme)
}
