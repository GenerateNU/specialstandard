package student

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetStudentSessions(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	// Validate that ID is a valid UUID - fail fast
	parsedID, err := uuid.Parse(studentID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID")
	}

	sessions, err := h.studentRepository.GetStudentSessions(c.Context(), parsedID)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get student sessions", "id", studentID, "err", err)
		return errs.InternalServerError("Failed to retrieve student sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
