package student

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetStudentRatings(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	// Validate that ID is a valid UUID - fail fast
	parsedID, err := uuid.Parse(studentID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID")
	}

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Parse filter parameters
	filter := &models.GetStudentSessionsRatingsRequest{}
	if err := c.QueryParser(filter); err != nil {
		slog.Error("Query parsing failed", "error", err, "query", c.OriginalURL())
		return errs.BadRequest("Error parsing filter parameters.")
	}

	if validationErrors := xvalidator.Validator.Validate(filter); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Validate that start date is before end date if both are provided
	if filter.StartDate != nil && filter.EndDate != nil {
		if filter.StartDate.After(*filter.EndDate) {
			return errs.BadRequest("Start date must be before end date")
		}
	}

	sessions, err := h.studentRepository.GetStudentRatings(c.Context(), parsedID, pagination, filter)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get student sessions", "id", studentID, "err", err)
		return errs.InternalServerError("Failed to retrieve student sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
