package therapist

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This is our function to send a request to our GetTgeraoustByID function
// We check for pretty errors in this function
func (h *Handler) GetTherapistByID(c *fiber.Ctx) error {
	therapistID := c.Params("id")

	_, err := uuid.Parse(therapistID)

	// Check if UUID is valid
	if err != nil {
		return errs.NotFound("The givenID was not found: ", therapistID)
	}

	// Checking for no ID given
	if therapistID == "" {
		return errs.NotFound("Given Empty ID")
	}

	therapist, err := h.therapistRepository.GetTherapistByID(c.Context(), therapistID)

	// Here we parse the bad request which is recieved
	if err != nil {
		slog.Error("Error updating document", "error", err)
		return errs.BadRequest("Failed to parse given request")
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}
