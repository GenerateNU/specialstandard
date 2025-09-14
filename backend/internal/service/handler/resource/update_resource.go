package resource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) UpdateResource(c *fiber.Ctx) error {
	idStr := c.Params("id")
	resourceId, err := uuid.Parse(idStr)
	if err != nil {
		return errs.InvalidRequestData(map[string]string{"id": "invalid UUID"})
	}

	slog.Info("Raw request body", "body", string(c.Body()))

	var resourceBody models.UpdateResourceBody
	if c.BodyParser(&resourceBody) != nil {
		return errs.InvalidRequestData(map[string]string{"body": "invalid body"})
	}

	res, err := h.resourceRepository.UpdateResource(c.Context(), resourceId, resourceBody)
	if err != nil {
		return errs.InternalServerError()
	}
	return c.Status(fiber.StatusOK).JSON(res)
}
