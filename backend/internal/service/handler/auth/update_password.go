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

	fmt.Printf("Received update password request\n")

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("Error parsing request body: %v\n", err)
		return errs.BadRequest("Invalid request body")
	}

	fmt.Printf("New password: %s\n", payload.Password)

	if payload.Password == "" {
		fmt.Println("New password is required")
		return errs.BadRequest("New password is required")
	}

	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Reset token is missing"})
	}

	fmt.Printf("Reset token: %s\n", token)

	err := auth.SupabaseUpdatePassword(&h.config, token, payload.Password)
	if err != nil {
		fmt.Printf("Password update failed: %v\n", err)
		return errs.InternalServerError(fmt.Sprintf("Password update failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
