package school

import (
	"specialstandard/internal/errs"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// GetSchools handles GET /schools
func (h *Handler) GetSchools(c *fiber.Ctx) error {
	// Check if district_id query param is provided
	districtIDStr := c.Query("district_id")

	if districtIDStr != "" {
		districtID, err := strconv.Atoi(districtIDStr)
		if err != nil {
			return errs.BadRequest("Invalid district ID")
		}
		
		schools, err := h.schoolRepository.GetSchoolsByDistrict(c.Context(), districtID)
		if err != nil {
			return errs.InternalServerError("Failed to fetch schools")
		}
		return c.Status(fiber.StatusOK).JSON(schools)
	}
	
	// Get all schools
	schools, err := h.schoolRepository.GetSchools(c.Context())
	if err != nil {
		return errs.InternalServerError("Failed to fetch schools")
	}
	
	return c.Status(fiber.StatusOK).JSON(schools)
}