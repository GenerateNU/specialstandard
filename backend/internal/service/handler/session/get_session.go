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

		// Check for duplicates if we have any IDs, and uses that one cancellation property that iforgot what what it was called
		if len(uuidStudentIDs) > 0 && checkForDuplicates(uuidStudentIDs) {
			return errs.BadRequest("Given multiple of the same students")
		}
	}

	// we do not need to check for invalid uuid in request body
	// because since it is a list of UUIDs, if they provide a non-uuid it will auto-error

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

	repoFilter := &models.GetSessionRepositoryRequest{
		StartTime:  filter.StartTime,
		EndTime:    filter.EndTime,
		Month:      filter.Month,
		Year:       filter.Year,
		StudentIDs: nil,
	}
	
	// Only set StudentIDs if we have valid UUIDs
	if len(uuidStudentIDs) > 0 {
		repoFilter.StudentIDs = &uuidStudentIDs
	}


	sessions, err := h.sessionRepository.GetSessions(c.Context(), pagination, repoFilter)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session", "err", err)
		return errs.InternalServerError("Failed to retrieve sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}

// Basic function to check for duplicates in a list of uuids!
func checkForDuplicates(ids []uuid.UUID) bool {
	seen := make(map[uuid.UUID]bool)
	for _, i := range ids {
		if seen[i] {
			return true
		}
		seen[i] = true
	}
	return false
}
