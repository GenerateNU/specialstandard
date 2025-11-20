package auth

import (
	"fmt"
	"os"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) ForgotPassword(c *fiber.Ctx) error {
	var payload struct {
		Email string `json:"email"`
	}

	fmt.Printf("Received forgot password request for email: %s\n", payload.Email)

	if err := c.BodyParser(&payload); err != nil {
		fmt.Printf("Error parsing request body: %v\n", err)
		return errs.BadRequest("Invalid request body")
	}

	if payload.Email == "" {
		fmt.Println("Email is required")
		return errs.BadRequest("Email is required")
	}

	// Get the frontend URL from environment or use default
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		// Fallback for local development
		frontendURL = "http://localhost:3000"
	}

	// Construct the redirect URL for the password reset link
	redirectURL := frontendURL + "/auth/reset-password"

	err := auth.SupabaseForgotPassword(&h.config, payload.Email, redirectURL)
	if err != nil {
		fmt.Printf("Password reset request failed: %v\n", err)
		return errs.InternalServerError(fmt.Sprintf("Password reset request failed: %v", err))
	}

	return c.SendStatus(fiber.StatusOK)
}
