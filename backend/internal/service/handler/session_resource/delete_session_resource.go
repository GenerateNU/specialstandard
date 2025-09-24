package session_resource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteSessionResource(c *fiber.Ctx) error {
	var sessionResource models.DeleteSessionResource

	if err := c.BodyParser(&sessionResource); err != nil {
		return errs.InvalidJSON("Failed to parse DeleteSessionResource data")
	}

	if validationErrors := h.validator.Validate(sessionResource); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	err := h.sessionResourceRepository.DeleteSessionResource(c.Context(), sessionResource)
	if err != nil {
		slog.Error("Failed to delete session_resource", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "23503") || strings.Contains(errStr, "relationship not found"): // foreign key violation
			return errs.NotFound("session or resource not found")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Delete Session Resource")
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
