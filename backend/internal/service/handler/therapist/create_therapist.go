package therapist

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTherapist(c *fiber.Ctx) error {
	var therapist models.CreateTherapistInput

	if err := c.BodyParser(&therapist); err != nil {
		return errs.InvalidJSON("Failed to parse therapist data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(therapist); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	createdTherapist, err := h.therapistRepository.CreateTherapist(c.Context(), &therapist)
	if err != nil {
		// Specific error handling with custom messages
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid reference to related data")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database connection error")
		default:
			return errs.InternalServerError(errStr)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(createdTherapist)
}
