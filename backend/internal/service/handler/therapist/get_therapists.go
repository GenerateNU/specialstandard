package therapist

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetTherapists(c *fiber.Ctx) error {
	sessions, err := h.therapistRepository.GetTherapists(c.Context())

	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
