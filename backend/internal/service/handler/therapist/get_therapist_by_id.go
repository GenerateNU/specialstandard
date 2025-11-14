package therapist

import (
	"log/slog"
	"specialstandard/internal/errs"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This is our function to send a request to our GetTherapistByID function
func (h *Handler) GetTherapistByID(c *fiber.Ctx) error {
	therapistID := c.Params("id")

	_, err := uuid.Parse(therapistID)

	// Check if UUID is valid
	if err != nil {
		return errs.BadRequest("Invalid UUID format for ID : ", therapistID)
	}

	therapist, err := h.therapistRepository.GetTherapistByID(c.Context(), therapistID)

	if err != nil {
		slog.Error("Failed to get therapist by ID", "therapist_id", therapistID, "err", err)
		errStr := strings.ToLower(err.Error())
		switch {
		case strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "no rows in result set"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Therapist not found",
			})
		case strings.Contains(errStr, "connection refused"):
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error": "Service unavailable",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}
