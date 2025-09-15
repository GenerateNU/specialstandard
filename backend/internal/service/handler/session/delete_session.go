package session

import (
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
		return errs.InternalServerError("Internal Server Error")
	} else if message == "" {
		return errs.NotFound("Session Not Found")
	}

	return c.Status(fiber.StatusOK).JSON(message)
}
