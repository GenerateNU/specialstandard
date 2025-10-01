package student

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetStudentsQuery represents the query parameters for filtering students
type GetStudentsQuery struct {
	Grade       string `query:"grade" validate:"omitempty"`
	TherapistID string `query:"therapist_id" validate:"omitempty,uuid"`
	Name        string `query:"name" validate:"omitempty"`
	utils.Pagination
}

func (h *Handler) GetStudents(c *fiber.Ctx) error {
	var query GetStudentsQuery

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

	// Validate query parameters
	if validationErrors := xvalidator.Validator.Validate(query); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
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

	// Call repository with extracted parameters
	students, err := h.studentRepository.GetStudents(
		c.Context(),
		query.Grade,
		therapistID,
		query.Name,
		query.Pagination,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(students)
}
