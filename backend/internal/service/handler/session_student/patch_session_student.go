package sessionstudent

import (
	"log/slog"
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// PatchSessionStudent handler
func (h *Handler) PatchSessionStudent(c *fiber.Ctx) error {
	var sessionStudent models.PatchSessionStudentInput

	if err := c.BodyParser(&sessionStudent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate required fields
	if sessionStudent.SessionID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}
	if sessionStudent.StudentID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	updatedSessionStudent, err := h.sessionStudentRepository.PatchSessionStudent(c.Context(), &sessionStudent)
	if err != nil {
		slog.Error("Failed to patch session student", "session_id", sessionStudent.SessionID, "student_id", sessionStudent.StudentID, "err", err)

		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "no rows affected") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "no rows in result set"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session student relationship not found",
			})
		case strings.Contains(errStr, "foreign key"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Reference",
			})
		case strings.Contains(errStr, "check constraint"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Violated a check constraint",
			})
		case strings.Contains(errStr, "connection refused"):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database Connection Error",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to Update SessionStudent",
			})
		}
	}
	return c.Status(fiber.StatusOK).JSON(updatedSessionStudent)
}
