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

	token := c.Query("token", "missing")

	switch token {
	case "missing":
		token = c.Cookies("jwt", "")
	case "":
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reset token is missing"})
	}

	err := auth.SupabaseUpdatePassword(&h.config, token, payload.Password)
	if err != nil {
		return errs.InternalServerError(fmt.Sprintf("Password update failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
