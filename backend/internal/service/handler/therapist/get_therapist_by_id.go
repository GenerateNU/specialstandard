package therapist

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetTherapistByID(c *fiber.Ctx) error {
	therapistID := c.Params("id")

	therapist, err := h.therapistRepository.GetTherapistByID(c.Context(), therapistID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}

