package auth

import (
	"fmt"
	"log/slog"
	"os"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var cred Credentials

	if err := c.BodyParser(&cred); err != nil {
		return errs.BadRequest(fmt.Sprintf("Invalid Request Body: %v", cred))
	}

	signInResponse, err := auth.SupabaseLogin(&h.config, cred.Email, cred.Password)
	if err != nil {
		slog.Error("Supabase Login Error: ", "err", err.Error())

		// Extract the actual message from HTTPError
		if httpErr, ok := err.(errs.HTTPError); ok {
			return errs.Unauthorized(httpErr.Message.(string))
		}

		return errs.Unauthorized("Invalid credentials")
	}

	fmt.Println(cred.RememberMe)

	var cookieExp time.Time
	if cred.RememberMe {
		cookieExp = time.Now().Add(7 * 24 * time.Hour)
	} else {
		cookieExp = time.Time{}
	}

	// Check if running in production
	isProduction := os.Getenv("ENV") == "production"

	c.Cookie(&fiber.Cookie{
		Name:     "userID",
		Value:    signInResponse.User.ID.String(),
		Expires:  cookieExp,
		Secure:   isProduction,
		SameSite: "None",
		Path:     "/",
		Domain:   "",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    signInResponse.AccessToken,
		Expires:  cookieExp,
		Secure:   isProduction,
		HTTPOnly: true,   // Recommended for JWT security
		SameSite: "None", // Changed from "Lax" to "None" for cross-origin
		Path:     "/",
		Domain:   "", // Leave empty or set to specific domain
	})

	return c.Status(fiber.StatusOK).JSON(signInResponse)
}
