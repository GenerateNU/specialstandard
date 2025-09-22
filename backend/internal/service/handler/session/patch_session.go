package session

import (
	"errors"
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) PatchSessions(c *fiber.Ctx) error {
	var session models.PatchSessionInput

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.BadRequest("Parsing Error with Invalid ID Format. ID: " + id.String())
	}

	if err := c.BodyParser(&session); err != nil {
		return errs.InvalidJSON("Failed to parse PatchSessionInput data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(session); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	updatedSession, err := h.sessionRepository.PatchSession(c.Context(), id, &session)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("Session Not Found")
		}

		slog.Error("Failed to patch session", "id", id, "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Update Session")
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedSession)
}
