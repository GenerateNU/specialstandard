package session

import (
	"errors"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetSessionByID(c *fiber.Ctx) error {
	sessionID := c.Params("id")

	// Checking for no ID given
	if sessionID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	// Validate that ID is a valid UUID - fail fast
	_, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.BadRequest("Invalid UUID format")
	}

	session, err := h.sessionRepository.GetSessionByID(c.Context(), sessionID)
	if err != nil {
		// Check if it's a "no rows found" error using pgx's error constant
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("Session not found")
		}
		// For all other database errors, return internal server error without exposing details
		return errs.InternalServerError("Failed to retrieve session")
	}

	return c.Status(fiber.StatusOK).JSON(session)
}
