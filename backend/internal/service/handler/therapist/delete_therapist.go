package therapist

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// This is our function to send a request to our GetTgeraoustByID function
// We check for pretty errors in this function
func (h *Handler) DeleteTherapist(c *fiber.Ctx) error {
	therapistID := c.Params("id")

	_, err := uuid.Parse(therapistID)

	// Check if UUID is valid
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	err = h.therapistRepository.DeleteTherapist(c.Context(), therapistID)

	// Here we parse the bad request which is recieved
	if err != nil {
		slog.Error("Error deleting document", "error", err)
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Therapist deleted successfully",
	})
}
