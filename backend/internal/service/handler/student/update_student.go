package student

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"specialstandard/internal/models"
	"strings"
	"time"
)

func (h *Handler) UpdateStudent(c *fiber.Ctx) error {
	// Get ID from URL parameter
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)

	// Check if UUID is valid
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	var req models.UpdateStudentInput
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Get existing student for merging (don't check for "not found" errors here)
	existingStudent, err := h.studentRepository.GetStudent(c.Context(), id)
	if err != nil {
		// Generic database error - let UpdateStudent handle "not found" case
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
		if *req.DOB == "" {
			// Empty string means set to NULL
			existingStudent.DOB = nil
		} else {
			dob, err := time.Parse("2006-01-02", *req.DOB)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid date format. Use YYYY-MM-DD",
				})
			}
			existingStudent.DOB = &dob
		}
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
		if *req.Grade == "" {
			// Empty string means set to NULL
			existingStudent.Grade = nil
		} else {
			existingStudent.Grade = req.Grade
		}
	}
	if req.IEP != nil {
		if *req.IEP == "" {
			// Empty string means set to NULL
			existingStudent.IEP = nil
		} else {
			existingStudent.IEP = req.IEP
		}
	}

	// Save updated student - let this call handle "student not found" errors
	updatedStudent, err := h.studentRepository.UpdateStudent(c.Context(), existingStudent)
	if err != nil {
		// Check if student was not found during update
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