package session_resource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSessionResources(c *fiber.Ctx) error {
	sessionId, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.InvalidRequestData(map[string]string{"id": "Invalid UUID format"})
	}

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	var resources []models.Resource
	if resources, err = h.sessionResourceRepository.GetResourcesBySessionID(c.Context(), sessionId, pagination); err != nil {
		slog.Error("Failed to get session's resources", "id", sessionId, "err", err)
		return errs.InternalServerError("Failed to retrieve session resources", err.Error())
	}

	if resources == nil {
		resources = []models.Resource{}
	}

	return c.JSON(resources)
}
