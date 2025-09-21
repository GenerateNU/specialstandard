package session

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSessionStudents(c *fiber.Ctx) error {
	sessionID := c.Params("id")

	if sessionID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	// Validate that ID is a valid UUID - fail fast
	parsedID, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID")
	}

	students, err := h.sessionRepository.GetSessionStudents(c.Context(), parsedID)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session students", "id", sessionID, "err", err)
		return errs.InternalServerError("Failed to retrieve session students")
	}

	return c.Status(fiber.StatusOK).JSON(students)
}
