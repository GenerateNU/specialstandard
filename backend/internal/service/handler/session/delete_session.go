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
		return errs.BadRequest("ID: Parsing Error")
	}

	message, err := h.sessionRepository.DeleteSessions(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete session", "id", id, "err", err)
		return errs.InternalServerError("Internal Server Error")
	} else if message == "not found" {
		return errs.NotFound("Session Not Found")
	}

	return c.Status(fiber.StatusOK).JSON(message)
}
