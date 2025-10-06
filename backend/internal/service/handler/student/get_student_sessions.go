package student

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"
	"time"

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

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Parse filter parameters - use QueryParser for standard types, manual for custom date formats
	filter := &models.GetStudentSessionsRequest{}
	if err := c.QueryParser(filter); err != nil {
		slog.Error("Query parsing failed", "error", err, "query", c.OriginalURL())
		return errs.BadRequest("Error parsing filter parameters.")
	}

	if validationErrors := xvalidator.Validator.Validate(filter); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Convert to repository request - start with standard fields from QueryParser
	repositoryFilter := &models.GetStudentSessionsRepositoryRequest{
		Month:   filter.Month,
		Year:    filter.Year,
		Present: filter.Present,
	}

	// Manually parse date fields since QueryParser expects RFC3339 format but we want YYYY-MM-DD
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			repositoryFilter.StartDate = &startDate
		} else {
			return errs.BadRequest("Invalid startDate format. Use YYYY-MM-DD")
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			repositoryFilter.EndDate = &endDate
		} else {
			return errs.BadRequest("Invalid endDate format. Use YYYY-MM-DD")
		}
	}

	// Only pass filter if there are actual filter parameters
	var finalFilter *models.GetStudentSessionsRepositoryRequest
	if repositoryFilter.StartDate != nil || repositoryFilter.EndDate != nil || 
	   repositoryFilter.Month != nil || repositoryFilter.Year != nil || repositoryFilter.Present != nil {
		finalFilter = repositoryFilter
	}

	sessions, err := h.studentRepository.GetStudentSessions(c.Context(), parsedID, pagination, finalFilter)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get student sessions", "id", studentID, "err", err)
		return errs.InternalServerError("Failed to retrieve student sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
