package student

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetStudents(c *fiber.Ctx) error {
	students, err := h.studentRepository.GetStudents(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(students)
}

