package therapist

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This is our function to send a request to our PatchTherapist function
// We check for pretty errors in this function
func (h *Handler) PatchTherapist(c *fiber.Ctx) error {
	therapistID := c.Params("id")
	// Ensure id is valid!
	_, err := uuid.Parse(therapistID)

	// Check if UUID is valid
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}
	var updatedValue models.UpdateTherapist

	if err := c.BodyParser(&updatedValue); err != nil {
		return errs.InvalidJSON("Failed to parse therapist data")
	}

	// Validate using XValidator
	if validationErrors := h.validator.Validate(updatedValue); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	therapist, err := h.therapistRepository.PatchTherapist(c.Context(), therapistID, &updatedValue)

	// Here we parse the bad request which is recieved
	if err != nil {
		// Here we parse the bad request which is recieved
		slog.Error("Failed to patch therapist", "therapist_id", therapistID, "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "no rows affected") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "no rows in result set"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Therapist not found",
			})
		case strings.Contains(errStr, "foreign key"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Reference",
			})
		case strings.Contains(errStr, "check constraint"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Violated a check constraint",
			})
		case strings.Contains(errStr, "connection refused"):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database Connection Error",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to Update Therapist",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}
