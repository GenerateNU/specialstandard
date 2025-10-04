package student

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PromoteStudents(c *fiber.Ctx) error {
	var promoteStudents models.PromoteStudentsInput

	if err := c.BodyParser(&promoteStudents); err != nil {
		return errs.InvalidJSON("Failed to parse PromoteStudentsInput")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(promoteStudents); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	err := h.studentRepository.PromoteStudents(c.Context(), promoteStudents)
	if err != nil {
		slog.Error("Failed to promote students", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Therapist/Student Reference")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Promote Students")
		}
	}

	return c.SendStatus(fiber.StatusNoContent)
}
