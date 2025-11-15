package auth

import (
	"fmt"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var payload struct {
		Email string `json:"email"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	if payload.Email == "" {
		return errs.BadRequest("Email is required")
	}

	err := auth.SupabaseForgotPassword(&h.config, payload.Email)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Password reset request failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
