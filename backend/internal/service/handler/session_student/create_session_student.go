package sessionstudent

import (
	"specialstandard/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateSessionStudent(c *fiber.Ctx) error {
	var req models.CreateSessionStudentInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	if req.SessionID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}
	if req.StudentID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	sessionStudent, err := h.sessionStudentRepository.CreateSessionStudent(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session student",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sessionStudent)
}
