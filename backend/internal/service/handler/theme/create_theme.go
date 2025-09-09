package theme

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTheme(c *fiber.Ctx) error {
	var theme models.CreateThemeInput

	if err := c.BodyParser(&theme); err != nil {
		return errs.InvalidJSON("Failed to parse theme data")
	}

	// Validate input
	validationErrors := make(map[string]string)

	if theme.Name == "" {
		validationErrors["name"] = "Name is required"
	}

	if theme.Month < 1 || theme.Month > 12 {
		validationErrors["month"] = "Month must be between 1 and 12"
	}

	if theme.Year < 1900 || theme.Year > 2100 {
		validationErrors["year"] = "Year must be between 1900 and 2100"
	}

	if len(validationErrors) > 0 {
		return errs.InvalidRequestData(validationErrors)
	}

	createdTheme, err := h.themeRepository.CreateTheme(c.Context(), &theme)
	if err != nil {
		// Specific error handling with custom messages
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "duplicate key"):
			return errs.Conflict("Theme with this name already exists for the specified month and year")
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
