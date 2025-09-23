package sessionstudent

import (
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DeleteSessionStudent handler
func (h *Handler) DeleteSessionStudent(c *fiber.Ctx) error {
	var req models.DeleteSessionStudentInput

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate required fields
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

	err := h.sessionStudentRepository.DeleteSessionStudent(c.Context(), &req)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session student relationship not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete session student",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
