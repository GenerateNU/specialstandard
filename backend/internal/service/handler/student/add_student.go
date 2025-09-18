package student

import (
	"github.com/gofiber/fiber/v2"
	"specialstandard/internal/models"
	"time"
	"github.com/google/uuid"
)

func (h *Handler) AddStudent(c *fiber.Ctx) error {
	var req struct {
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		DOB         string `json:"dob"`
		TherapistID string `json:"therapist_id"`
		Grade       string `json:"grade"`
		IEP         string `json:"iep"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "Invalid JSON format",
    })
	}
	
	// Parse date
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid date format. Use YYYY-MM-DD",
		})
	}
	
	therapistID, err := uuid.Parse(req.TherapistID)
	// Check if UUID is valid
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid therapist ID format",
		})
	}

	// // TODO: UNCOMMENT Validate therapist exists using GetTherapist
	// _, err = h.Repository.Therapist.GetTherapist(c.Context(), therapistID)
	// if err != nil {
	// 	// Check for "not found" error patterns
	// 	if strings.Contains(err.Error(), "no rows") || err == sql.ErrNoRows {
	// 		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
	// 			"error": "Therapist not found",
	// 		})
	// 	}
	// 	// Other database errors
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Database error while validating therapist",
	// 	})
	// }
	
	student := models.Student{
		ID:          uuid.New(), 
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DOB:         dob,
		TherapistID: therapistID,
		Grade:       req.Grade,
		IEP:         req.IEP,
	}
	
	student,  err = h.studentRepository.AddStudent(c.Context(), student)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}
	
	return c.Status(fiber.StatusCreated).JSON(student)
}