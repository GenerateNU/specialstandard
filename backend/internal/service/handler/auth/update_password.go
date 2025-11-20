package auth

import (
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdatePassword(c *fiber.Ctx) error {
	var payload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OTP      string `json:"otp"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	if payload.Email == "" {
		return errs.BadRequest("Email is required")
	}

	if payload.Password == "" {
		return errs.BadRequest("Password is required")
	}

	if payload.OTP == "" {
		return errs.BadRequest("OTP is required")
	}

	// Verify OTP and update password
	err := auth.VerifyOTPAndResetPassword(&h.config, payload.Email, payload.OTP, payload.Password)
	if err != nil {
		return errs.BadRequest(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
