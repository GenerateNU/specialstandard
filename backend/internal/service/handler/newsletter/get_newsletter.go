package newsletter

import (
	"log/slog"
	"specialstandard/internal/errs"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GET /api/v1/newsletter/by-date?date=YYYY-MM-DD
func (h *Handler) GetNewsletterByDate(c *fiber.Ctx) error {
	dateStr := c.Query("date")
	if dateStr == "" {
		return errs.BadRequest("Missing date query parameter")
	}

	// Validate date format
	parsedDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errs.BadRequest("Invalid date format, use YYYY-MM-DD")
	}

	newsletter, err := h.Repo.GetNewsletterByDate(c.Context(), parsedDate)
	if err != nil || newsletter == nil {
		return errs.NotFound("newsletter", "No newsletter found for this date")
	}

	// Extract the object key from S3URL
	// Handle various formats: "s3://bucket/path", "/path", or just "path"
	key := ""
	if newsletter.S3URL != "" {
		key = newsletter.S3URL

		// Remove s3:// prefix and bucket name if present
		if strings.HasPrefix(key, "s3://") {
			key = strings.TrimPrefix(key, "s3://")
			// Remove bucket name
			parts := strings.SplitN(key, "/", 2)
			if len(parts) == 2 {
				key = parts[1]
			}
		}

		// Remove leading slash (same as resource handler does)
		key = strings.TrimPrefix(key, "/")
	}

	if key == "" {
		return errs.InternalServerError("Newsletter S3 key is empty")
	}

	// Log for debugging
	slog.Info("Newsletter: About to generate presigned URL",
		"s3_url_raw", newsletter.S3URL,
		"key_extracted", key,
		"key_length", len(key),
	)

	// Generate presigned URL (same as resource handler)
	url, err := h.s3Client.GeneratePresignedURL(c.Context(), key, time.Hour)
	if err != nil {
		slog.Error("Failed to generate presigned URL for newsletter",
			"key", key,
			"error", err,
		)
		return errs.InternalServerError("Failed to generate presigned URL")
	}

	slog.Info("Newsletter: Generated presigned URL",
		"key", key,
		"url_length", len(url),
		"url_prefix", url[:100], // First 100 chars
	)

	return c.JSON(fiber.Map{"s3_url": url})
}
