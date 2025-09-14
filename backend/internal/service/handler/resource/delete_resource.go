package resource

import (
	"log/slog"
	"specialstandard/internal/errs"

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
	slog.Info("deleting resource", "id", resourceId)
	if err != nil {
		slog.Error("failed to delete resource", "error", err)
		return err
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
