package student

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"
	"strconv"
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

	// Manually parse filter parameters
	filter := &models.GetStudentSessionsRepositoryRequest{}

	// Parse date parameters
	if startDateStr := c.Query("startDate"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		} else {
			return errs.BadRequest("Invalid startDate format. Use YYYY-MM-DD")
		}
	}

	if endDateStr := c.Query("endDate"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			filter.EndDate = &endDate
		} else {
			return errs.BadRequest("Invalid endDate format. Use YYYY-MM-DD")
		}
	}

	// Parse month parameter
	if monthStr := c.Query("month"); monthStr != "" {
		if month, err := strconv.Atoi(monthStr); err == nil {
			if month < 1 || month > 12 {
				return errs.BadRequest("Month must be between 1 and 12")
			}
			filter.Month = &month
		} else {
			return errs.BadRequest("Invalid month format")
		}
	}

	// Parse year parameter
	if yearStr := c.Query("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			if year < 1776 || year > 2200 {
				return errs.BadRequest("Year must be between 1776 and 2200")
			}
			filter.Year = &year
		} else {
			return errs.BadRequest("Invalid year format")
		}
	}

	// Parse present parameter
	if presentStr := c.Query("present"); presentStr != "" {
		if presentStr == "true" {
			present := true
			filter.Present = &present
		} else if presentStr == "false" {
			present := false
			filter.Present = &present
		} else {
			return errs.BadRequest("Present must be 'true' or 'false'")
		}
	}

	// Only pass filter if there are actual filter parameters
	var repositoryFilter *models.GetStudentSessionsRepositoryRequest
	if filter.StartDate != nil || filter.EndDate != nil || filter.Month != nil || filter.Year != nil || filter.Present != nil {
		repositoryFilter = filter
	}

	sessions, err := h.studentRepository.GetStudentSessions(c.Context(), parsedID, pagination, repositoryFilter)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get student sessions", "id", studentID, "err", err)
		return errs.InternalServerError("Failed to retrieve student sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
