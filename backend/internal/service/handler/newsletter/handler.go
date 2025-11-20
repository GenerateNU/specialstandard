package newsletter

import (
	"net/http"
	"time"

	"specialstandard/internal/storage"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	Repo storage.NewsletterRepository
}

func NewHandler(repo storage.NewsletterRepository) *Handler {
	return &Handler{Repo: repo}
}

// GET /api/v1/newsletter/by-date?date=YYYY-MM-DD
func (h *Handler) GetNewsletterByDate(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing date query parameter"})
	}
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid date format, use YYYY-MM-DD"})
	}
	newsletter, err := h.Repo.GetNewsletterByDate(c.Context(), parsedDate)
	if err != nil || newsletter == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "No newsletter found for this date"})
	}
	return c.JSON(fiber.Map{"s3_url": newsletter.S3URL})
}
