package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type GetSessionStudentsQuery struct {
	TherapistID string `query:"therapist_id"`
	utils.Pagination
}

func (h *Handler) GetSessionStudents(c *fiber.Ctx) error {
	sessionID := c.Params("id")

	if sessionID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	var query GetSessionStudentsQuery

	// Set default pagination first
	query.Pagination = utils.NewPagination()

	if err := c.QueryParser(&query); err != nil {
		return errs.BadRequest("Invalid Query Parameters")
	}

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	if c.Query("therapist_id") == "" {
		query.TherapistID = ""
	}

	// Validate therapist_id only if non-empty
	if query.TherapistID != "" {
		if _, err := uuid.Parse(query.TherapistID); err != nil {
			return errs.BadRequest("Invalid therapist_id format")
		}
	}

	if validationErrors := xvalidator.Validator.Validate(query); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Validate that session ID is a valid UUID - fail fast
	parsedID, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID")
	}

	// Convert therapist_id string to UUID if provided
	var therapistID uuid.UUID
	if query.TherapistID != "" {
		parsedUUID, err := uuid.Parse(query.TherapistID)
		if err != nil {
			return errs.BadRequest("Invalid therapist_id format")
		}
		therapistID = parsedUUID
	}

	students, err := h.sessionRepository.GetSessionStudents(c.Context(), parsedID, query.Pagination, therapistID)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session students", "id", sessionID, "err", err)
		return errs.InternalServerError("Failed to retrieve session students")
	}

	return c.Status(fiber.StatusOK).JSON(students)
}