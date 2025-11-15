package district

import (
	"fmt"
	"specialstandard/internal/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetDistrictByID handles GET /districts/:id
func (h *Handler) GetDistrictByID(c *fiber.Ctx) error {
	sessionID := c.Params("id")

	// Checking for no ID given
	if sessionID == "" {
		return errs.BadRequest("Given Empty ID")
	}

	// Validate that ID is a valid UUID - fail fast
	_, err := uuid.Parse(sessionID)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("Invalid UUID format for ID '%s'", sessionID))
	}

	id, err := strconv.Atoi(sessionID)

	if err != nil {
		return errs.BadRequest("Invalid district ID")
	}
	
	district, err := h.districtRepository.GetDistrictByID(c.Context(), id)
	if err != nil {
		return errs.NotFound("District not found")
	}
	
	return c.Status(fiber.StatusOK).JSON(district)
}