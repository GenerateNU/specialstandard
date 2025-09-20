package theme

import (
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	var req models.UpdateThemeInput
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse theme data",
		})
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(req); len(validationErrors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": xvalidator.ConvertToMessages(validationErrors),
		})
	}

	// Update the theme
	updatedTheme, err := h.themeRepository.UpdateTheme(c.Context(), id, &req)
	if err != nil {
		// Check if theme was not found during update
		if strings.Contains(err.Error(), "no rows") || err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Theme not found",
			})
		}
		// Specific error handling with custom messages
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid reference to related data",
			})
		case strings.Contains(errStr, "connection refused"):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database connection error",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update theme",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedTheme)
}
