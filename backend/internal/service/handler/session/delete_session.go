package session

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteSessions(c *fiber.Ctx) error {
	id := c.Params("id")

	message, err := h.sessionRepository.DeleteSessions(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(message)
}
