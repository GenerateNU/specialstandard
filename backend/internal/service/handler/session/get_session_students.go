package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSessionStudents(c *fiber.Ctx) error {
	sessionID := c.Params("id")

	if sessionID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Validate that ID is a valid UUID - fail fast
	parsedID, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID")
	}

	students, err := h.sessionRepository.GetSessionStudents(c.Context(), parsedID, pagination)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session students", "id", sessionID, "err", err)
		return errs.InternalServerError("Failed to retrieve session students")
	}

	return c.Status(fiber.StatusOK).JSON(students)
}
