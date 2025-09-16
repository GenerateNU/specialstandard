package therapist

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

// This is our function to send a request to our GetTgeraoustByID function
// We check for pretty errors in this function
func (h *Handler) GetTherapistByID(c *fiber.Ctx) error {
	therapistID := c.Params("id")

	// Checking for no ID given
	if therapistID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	therapist, err := h.therapistRepository.GetTherapistByID(c.Context(), therapistID)

	// Here we parse the bad request which is recieved
	if err != nil {
		slog.Error("Error updating document", "error", errUpdate)
		return errs.BadRequest("There was an error parsing the given id!")
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}
