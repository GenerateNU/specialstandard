package sessionstudent

import (
	"log/slog"
	"specialstandard/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) PatchStudentSessionRatings(c *fiber.Ctx) error {
	var studentSessionRatings models.RateStudentSessionInput

	if err := c.BodyParser(&studentSessionRatings); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	// Validate required fields
	if studentSessionRatings.SessionID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}
	if studentSessionRatings.StudentID == (uuid.UUID{}) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student ID is required",
		})
	}

	// Validate rating enums
	validCategories := map[string]bool{
		"visual_cue":   true,
		"verbal_cue":   true,
		"gestural_cue": true,
		"engagement":   true,
	}

	validLevels := map[string]bool{
		"minimal":  true,
		"moderate": true,
		"maximal":  true,
		"low":      true,
		"high":     true,
	}

	for _, rating := range studentSessionRatings.Ratings {
		if !validCategories[rating.Category] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid rating category: " + rating.Category,
			})
		}
		if !validLevels[rating.Level] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid rating level: " + rating.Level,
			})
		}
	}

	student_session, ratings, err := h.sessionStudentRepository.RateStudentSession(c.Context(), &studentSessionRatings)
	if err != nil {
		slog.Error("Failed to rate session student", "session_id", studentSessionRatings.SessionID, "student_id", studentSessionRatings.StudentID, "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "no rows affected") ||
			strings.Contains(errStr, "not found") ||
			strings.Contains(errStr, "no rows in result set"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Session student relationship not found",
			})
		case strings.Contains(errStr, "foreign key"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid Reference",
			})
		case strings.Contains(errStr, "check constraint"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Violated a check constraint",
			})
		case strings.Contains(errStr, "connection refused"):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Database Connection Error",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to Update SessionStudent",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"sessionId": student_session.SessionID,
		"studentId": student_session.StudentID,
		"present":   student_session.Present,
		"notes":     student_session.Notes,
		"ratings":   ratings,
	})
}
