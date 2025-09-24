package session

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) DeleteSessions(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.BadRequest("Parsing Error: Invalid ID Format. ID: " + id.String())
	}

	err = h.sessionRepository.DeleteSession(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete session", "id", id, "err", err)
		return errs.InternalServerError("Internal Server Error")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Session deleted successfully",
	})
}
