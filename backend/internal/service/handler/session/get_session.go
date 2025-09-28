package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	var req models.GetSessionRequest
	if err := c.QueryParser(&req); err != nil {
		return errs.BadRequest("Error parsing request body.")
	}

	// check for no students in request body
	if len(req.StudentIDs) == 0 {
		return errs.BadRequest("No Student ID's recieved.")
	}



	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	sessions, err := h.sessionRepository.GetSessions(c.Context(), pagination)
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
