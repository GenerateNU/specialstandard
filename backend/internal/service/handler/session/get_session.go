package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	filter := &models.GetSessionRequest{}
	if err := c.QueryParser(filter); err != nil {
		slog.Error("Query parsing failed", "error", err, "query", c.OriginalURL())
		return errs.BadRequest("Error parsing request body.")
	}

	therapistID, err := uuid.Parse(filter.TherapistID)
	if err != nil {
		return errs.BadRequest("Invalid therapist ID format")
	}

	var uuidStudentIDs []uuid.UUID
	if filter.StudentIDs != nil && len(*filter.StudentIDs) > 0 {
		for _, idStr := range *filter.StudentIDs {
			if idStr == "" {
				continue // Skip empty strings
			}
			id, err := uuid.Parse(idStr)
			if err != nil {
				return errs.BadRequest("Invalid student ID format")
			}
			uuidStudentIDs = append(uuidStudentIDs, id)
		}
	}

	// check for valid time range in request body if time is given
	if filter.StartTime != nil && filter.EndTime != nil && filter.EndTime.Before(*filter.StartTime) {
		return errs.BadRequest("Given invalid time range.")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	if validationErrors := xvalidator.Validator.Validate(filter); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Only create repoFilter if there are actual filters to apply
	var repoFilter *models.GetSessionRepositoryRequest
	if filter.StartTime != nil || filter.EndTime != nil || filter.Month != nil || 
	   filter.Year != nil || len(uuidStudentIDs) > 0 {
		repoFilter = &models.GetSessionRepositoryRequest{
			StartTime:  filter.StartTime,
			EndTime:    filter.EndTime,
			Month:      filter.Month,
			Year:       filter.Year,
		}
		if len(uuidStudentIDs) > 0 {
			repoFilter.StudentIDs = &uuidStudentIDs
		}
	}

	sessions, err := h.sessionRepository.GetSessions(c.Context(), pagination, repoFilter, therapistID)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session", "err", err)
		return errs.InternalServerError("Failed to retrieve sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
