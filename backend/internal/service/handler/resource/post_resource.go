package resource

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"

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
		return errs.InternalServerError(err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(newResource)

}
