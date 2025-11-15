package sessionstudent

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetStudentAttendance(c *fiber.Ctx) error {
	studentID := c.Params("id")

	if studentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	// Validate that ID is a valid UUID
	parsedID, err := uuid.Parse(studentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is invalid UUID format",
		})
	}

	if parsedID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is invalid (nil)",
		})
	}

	dateFrom := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC) // old date as default - get all records
	if dateStr := c.Query("date_from"); dateStr != "" {
		var err error
		dateFrom, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date_from format",
			})
		}
	}

	dateTo := time.Now()
	if dateStr := c.Query("date_to"); dateStr != "" {
		var err error
		dateTo, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid date_to format",
			})
		}
	}

	presentCount, totalCount, err := h.sessionStudentRepository.GetStudentAttendance(c.Context(), models.GetStudentAttendanceParams{
		StudentID: parsedID,
		DateFrom:  dateFrom,
		DateTo:    dateTo,
	})
	if err != nil {
		slog.Error("Failed to get student attendance", "id", studentID, "err", err)
		return errs.InternalServerError("Failed to retrieve student sessions")
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"present_count": presentCount,
		"total_count":   totalCount,
	})
}
