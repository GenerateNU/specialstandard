package session_resource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PostSessionResource(c *fiber.Ctx) error {
	var sessionResource models.CreateSessionResource

	if err := c.BodyParser(&sessionResource); err != nil {
		return errs.InvalidJSON("Failed to parse PostSessionResource request body")
	}

	if validationErrors := h.validator.Validate(sessionResource); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	newSessionResource, err := h.sessionResourceRepository.PostSessionResource(c.Context(), sessionResource)
	if err != nil {
		slog.Error("Failed to post session_resource", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "23503"): // foreign key violation
			return errs.NotFound("session not found")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Create Session Resource")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(newSessionResource)
}
