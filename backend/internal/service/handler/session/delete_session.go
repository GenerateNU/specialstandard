package session

import (
	"specialstandard/internal/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteSessions(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return errs.BadRequest("Invalid ID: Parsing Error")
	}

	message, err := h.sessionRepository.DeleteSessions(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(message)
}
