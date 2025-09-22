package sessionstudent

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) DeleteSessionStudent(c *fiber.Ctx) error {
	var req struct {
		SessionID string `json:"session_id"`
		StudentID string `json:"student_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.SessionID == "" || req.StudentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID and Student ID are required",
		})
	}

	err := h.sessionStudentRepository.DeleteSessionStudent(c.Context(), req.SessionID, req.StudentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete session student",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
