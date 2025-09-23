package sessionresource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSessionResources(c *fiber.Ctx) (*[]models.Resource, error) {
	idStr := c.Params("id")
	if idStr == "" {
		return nil, errs.InvalidRequestData(map[string]string{"id": "Given Empty ID"})
	}

	sessionId, err := uuid.Parse(idStr)

	if err != nil {
		return nil, errs.InvalidRequestData(map[string]string{"id": "Invalid UUID format"})
	}

	var resources []models.Resource
	if resources, err = h.sessionResourceRepository.GetResourcesBySessionID(c.Context(), sessionId); err != nil {
		slog.Error("Failed to get session", "id", sessionId, "err", err)
		return nil, errs.InternalServerError("Failed to retrieve session resources", err.Error())
	}

	return &resources, nil
}
