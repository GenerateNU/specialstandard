package student

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"
	"time"
)

func (h *Handler) UpdateStudent(c *fiber.Ctx) error {
	// Get ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)

	if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "Invalid UUID format",
    })
}

	var req struct {
		FirstName   *string `json:"first_name,omitempty"`
		LastName    *string `json:"last_name,omitempty"`
		DOB         *string `json:"dob,omitempty"`
		TherapistID *string `json:"therapist_id,omitempty"`
		Grade       *string `json:"grade,omitempty"`
		IEP         *string `json:"iep,omitempty"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	existingStudent, err := h.studentRepository.GetStudent(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	// Update fields if provided
	if req.FirstName != nil {
		existingStudent.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existingStudent.LastName = *req.LastName
	}
	if req.DOB != nil {
		dob, err := time.Parse("2006-01-02", *req.DOB)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date format. Use YYYY-MM-DD",
			})
		}
		existingStudent.DOB = dob
	}
	if req.TherapistID != nil {
		therapistID, err := uuid.Parse(*req.TherapistID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid therapist ID format",
			})
		}
		existingStudent.TherapistID = therapistID
	}
	if req.Grade != nil {
		existingStudent.Grade = *req.Grade
	}
	if req.IEP != nil {
		existingStudent.IEP = *req.IEP
	}

	updatedStudent, err := h.studentRepository.UpdateStudent(c.Context(), existingStudent)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") || err.Error() == "sql: no rows in result set" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Student not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedStudent)

}