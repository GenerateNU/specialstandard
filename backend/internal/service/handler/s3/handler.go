package s3

import (
	"context"
	"net/http"
	"specialstandard/internal/s3_client"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	S3 *s3_client.Client
}

func NewHandler(s3Client *s3_client.Client) *Handler {
	return &Handler{S3: s3Client}
}

// POST /api/v1/s3/presign
func (h *Handler) GeneratePresignedURL(c *fiber.Ctx) error {
	type request struct {
		Key    string `json:"key"`
		Expiry int    `json:"expiry"`
	}
	var req request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	if req.Key == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Missing key"})
	}
	expiry := 900 // default 15 min
	if req.Expiry > 0 {
		expiry = req.Expiry
	}
	url, err := h.S3.GeneratePresignedURL(context.Background(), req.Key, time.Duration(expiry)*time.Second)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"url": url})
}
