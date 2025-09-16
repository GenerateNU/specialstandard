package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) PatchSessions(c *fiber.Ctx) error {
	var session models.PatchSessionInput

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.BadRequest("ID: Parsing Error")
	}

	if err := c.BodyParser(&session); err != nil {
		return errs.InvalidJSON("Failed to parse PatchSessionInput data")
	}

	updatedSession, err := h.sessionRepository.PatchSessions(c.Context(), id, &session)
	if err != nil {
		slog.Error("Failed to patch session", "id", id, "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "Not Found"):
			return errs.NotFound("Session Not Found")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Update Session")
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedSession)
}
