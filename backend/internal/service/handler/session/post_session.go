package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) PostSessions(c *fiber.Ctx) error {
	var session models.PostSessionInput

	// Parsing Session Inputs
	if err := c.BodyParser(&session); err != nil {
		return errs.InvalidJSON("Failed to parse PostSessionInput data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(session); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	var sessionIDs []uuid.UUID
	postSessionStudent := models.CreateSessionStudentInput{
		SessionIDs: sessionIDs,
		StudentIDs: nil,
		Present:    true,
		Notes:      nil,
	}
	if session.StudentIDs != nil {
		postSessionStudent.StudentIDs = *session.StudentIDs
	}

	hasStudents := session.StudentIDs != nil && len(*session.StudentIDs) > 0
	if hasStudents {
		for _, id := range *session.StudentIDs {
			if id == uuid.Nil {
				return errs.BadRequest("Student IDs must not contain empty UUIDs")
			}
		}
	}

	// Beginning Transaction
	tx, err := h.sessionRepository.GetDB().Begin(c.Context())
	if err != nil {
		return errs.InternalServerError("Failed to start transaction")
	}

	newSessions, err := h.sessionRepository.PostSession(c.Context(), tx, &session)
	if err != nil {
		slog.Error("Failed to post session", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Create Session")
		}
	}

	for _, newSession := range *newSessions {
		sessionIDs = append(sessionIDs, newSession.ID)
	}
	postSessionStudent.SessionIDs = sessionIDs

	if validationErrors := h.validator.Validate(postSessionStudent); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	_, err = h.sessionStudentRepository.CreateSessionStudent(c.Context(), tx, &postSessionStudent)
	if err != nil {
		rollbackErr := tx.Rollback(c.Context())
		if rollbackErr != nil {
			slog.Error("Rollback was not successful", "err", rollbackErr)
		}

		if strings.Contains(err.Error(), "unique_violation") || strings.Contains(err.Error(), "duplicate key") {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Student is already in this session",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create session student",
		})
	}

	err = tx.Commit(c.Context())
	if err != nil {
		return errs.InternalServerError("Failed to commit transaction")
	}

	return c.Status(fiber.StatusCreated).JSON(newSessions)
}
