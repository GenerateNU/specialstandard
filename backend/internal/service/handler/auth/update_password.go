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
	}

	if err := c.BodyParser(&payload); err != nil {
		return errs.BadRequest("Invalid request body")
	}

	if payload.Password == "" {
		return errs.BadRequest("New password is required")
	}

	// The token is typically passed as a query parameter from the reset link
	token := c.Query("token")
	if token == "" {
		return errs.BadRequest("Reset token is missing")
	}

	err := auth.SupabaseUpdatePassword(&h.config, token, payload.Password)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Password update failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
