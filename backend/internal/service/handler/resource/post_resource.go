package resource

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PostResource(c *fiber.Ctx) error {
	var resourceBody models.ResourceBody
	if err := c.BodyParser(&resourceBody); err != nil {
		return errs.InvalidRequestData(map[string]string{"body": "invalid body"})
	}

	if validationErrors := h.validator.Validate(&resourceBody); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	newResource, err := h.resourceRepository.CreateResource(c.Context(), resourceBody)
	if err != nil {
		// Check for foreign key constraint errors
		errStr := strings.ToLower(err.Error())
		if strings.Contains(errStr, "foreign key") ||
			strings.Contains(errStr, "violates foreign key constraint") ||
			strings.Contains(errStr, "invalid theme") ||
			strings.Contains(errStr, "23503") { // PostgreSQL foreign key violation code
			return errs.InvalidRequestData(map[string]string{
				"theme_id": "invalid theme",
			})
		}

		return errs.InternalServerError(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newResource)

}
