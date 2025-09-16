package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PostSessions(c *fiber.Ctx) error {
	var session models.PostSessionInput

	if err := c.BodyParser(&session); err != nil {
		return errs.InvalidJSON("Failed to parse PostSessionInput data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(session); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	newSession, err := h.sessionRepository.PostSessions(c.Context(), &session)
	if err != nil {
		slog.Error("Failed to delete session", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Create Session")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(newSession)
}
