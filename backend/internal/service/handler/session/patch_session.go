package session

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PatchSessions(c *fiber.Ctx) error {
	var session models.PatchSessionInput

	id := c.Params("id")

	if err := c.BodyParser(&session); err != nil {
		return errs.InvalidJSON("Failed to parse PatchSessionInput data")
	}

	updatedSession, err := h.sessionRepository.PatchSessions(c.Context(), id, &session)
	if err != nil {
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Update Session")
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedSession)
}
