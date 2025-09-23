package sessionstudent

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"

	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PatchSessionStudent(c *fiber.Ctx) error {
	var sessionStudent models.PatchSessionStudentInput

	if err := c.BodyParser(&sessionStudent); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate required fields
	if sessionStudent.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}
	if sessionStudent.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	updatedSessionStudent, err := h.sessionStudentRepository.PatchSessionStudent(c.Context(), sessionStudent.SessionID, sessionStudent.StudentID, &sessionStudent)
	if err != nil {
		if strings.Contains(err.Error(), "no rows affected") || strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session student relationship not found",
			})
		}

		slog.Error("Failed to patch session", "id", sessionStudent.SessionID, "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Update SessionStudent")
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedSessionStudent)
}
