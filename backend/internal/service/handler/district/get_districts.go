package district

import (
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

// GetDistricts handles GET /districts
func (h *Handler) GetDistricts(c *fiber.Ctx) error {
	districts, err := h.districtRepository.GetDistricts(c.Context())
	if err != nil {
		return errs.InternalServerError("Failed to fetch districts")
	}
	
	return c.Status(fiber.StatusOK).JSON(districts)
}