package auth

import (
	"fmt"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) UpdatePassword(c *fiber.Ctx) error {
	var payload struct {
		Password string `json:"password"`
		Token    string `json:"token"` // the recovery token from email
	}

	if err := c.BodyParser(&payload); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	if payload.Password == "" {
		return errs.BadRequest("New password is required")
	}

	// If token is provided, it's a password reset flow (unauthenticated)
	// If token is empty, it's from authenticated user (use JWT from cookie)
	token := payload.Token
	if token == "" {
		token = c.Cookies("jwt", "")
	}

	if token == "" {
		return errs.BadRequest("Authentication required")
	}

	err := auth.SupabaseUpdatePassword(&h.config, token, payload.Password)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Password update failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
