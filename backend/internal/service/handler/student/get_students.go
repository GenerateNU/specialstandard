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
	Grade       *int   `query:"grade" validate:"omitempty,oneof=-1 0 1 2 3 4 5 6 7 8 9 10 11 12"`
	TherapistID string `query:"therapist_id"`
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

	// Handle empty string query parameters - treat them as not provided
	// When ?grade= is present but empty, Fiber may parse it as a zero value
	if c.Query("grade") == "" {
		query.Grade = nil
	}
	if c.Query("name") == "" {
		query.Name = ""
	}
	if c.Query("therapist_id") == "" {
		query.TherapistID = ""
	}

	// Validate therapist_id only if non-empty (omitempty doesn't work with empty strings in query params)
	if query.TherapistID != "" {
		if _, err := uuid.Parse(query.TherapistID); err != nil {
			return errs.BadRequest("Invalid therapist_id format")
		}
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
		query.Grade, // Pass pointer directly - nil means no filter
		therapistID,
		query.Name,
		query.Pagination,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(students)
}
