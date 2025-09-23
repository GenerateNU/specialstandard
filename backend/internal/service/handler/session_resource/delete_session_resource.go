package sessionresource

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
		case strings.Contains(errStr, "foreign key"):
			if strings.Contains(errStr, "session_id") {
				return errs.NotFound("session not found")
			}
			if strings.Contains(errStr, "resource_id") {
				return errs.NotFound("resource not found")
			}
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Create Session Resource")
		}
	}

	return c.Status(fiber.StatusNoContent).JSON(nil)
}
