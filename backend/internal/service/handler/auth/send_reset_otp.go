package auth

import (
	"fmt"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SendResetOTP(c *fiber.Ctx) error {
	var payload struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	if payload.Email == "" {
		return errs.BadRequest("Email is required")
	}

	if payload.Token == "" {
		return errs.BadRequest("Reset token is required")
	}

	// Generate OTP and send email
	err := auth.SendResetOTP(&h.config, payload.Email)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Failed to send OTP: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
