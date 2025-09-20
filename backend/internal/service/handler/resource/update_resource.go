package resource

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) UpdateResource(c *fiber.Ctx) error {
	idStr := c.Params("id")
	resourceId, err := uuid.Parse(idStr)
	if err != nil {
		return errs.InvalidRequestData(map[string]string{"id": "invalid UUID"})
	}

	var resourceBody models.UpdateResourceBody
	if c.BodyParser(&resourceBody) != nil {
		return errs.InvalidRequestData(map[string]string{"body": "invalid body"})
	}

	if validationErrors := h.validator.Validate(&resourceBody); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	res, err := h.resourceRepository.UpdateResource(c.Context(), resourceId, resourceBody)

	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return errs.NotFound("resource", "not found")
		}
		return errs.InternalServerError(err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(res)
}
