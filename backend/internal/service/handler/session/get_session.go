package session

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	sessions, err := h.sessionRepository.GetSessions(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
