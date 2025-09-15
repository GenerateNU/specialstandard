package student

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
)

func (h *Handler) DeleteStudent(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	// Check if UUID is valid
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}
	
	if err := h.studentRepository.DeleteStudent(c.Context(), id); err != nil {
		// Student not found
		if strings.Contains(err.Error(), "no rows") || err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Student not found",
			})
		}
		// db error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	return c.SendStatus(fiber.StatusNoContent)
}