package sessionstudent

import (
	"specialstandard/internal/models"

	"strings"

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

	// Validate required fields
	for _, id := range req.SessionIDs {
		if id == (uuid.UUID{}) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Session ID is required",
			})
		}
	}

	for _, id := range req.StudentIDs {
		if id == (uuid.UUID{}) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Student ID is required",
			})
		}
	}

	db := h.sessionStudentRepository.GetDB()
	sessionStudents, err := h.sessionStudentRepository.CreateSessionStudent(c.Context(), db, &req)
	if err != nil {
		if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Student is already in this session",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session student",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(sessionStudents)
}
