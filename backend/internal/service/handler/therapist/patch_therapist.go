package therapist

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"

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

	// Checking for no ID given
	if therapistID == "" {
		return errs.BadRequest("Given Empty ID")
	}

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
		return errs.BadRequest("There was an error parsing the given id!")
	}

	return c.Status(fiber.StatusOK).JSON(therapist)
}
