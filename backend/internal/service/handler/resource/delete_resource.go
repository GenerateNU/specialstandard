package resource

import (
	"specialstandard/internal/errs"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) DeleteResource(c *fiber.Ctx) error {
	idStr := c.Params("id")
	resourceId, err := uuid.Parse(idStr)
	if err != nil {
		return errs.InvalidRequestData(map[string]string{"id": "invalid UUID"})
	}

	err = h.resourceRepository.DeleteResource(c.Context(), resourceId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return errs.NotFound("resource", "not found")
		}
		return errs.InternalServerError(err.Error())
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
