package auth

import (
	"fmt"
	"log/slog"
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
		slog.Error("Supabase Login Error: ", err)
		return errs.Unauthorized("Failed to login: ", err.Error())
	}

	fmt.Println(cred.RememberMe)

	var cookieExp time.Time
	if cred.RememberMe {
		cookieExp = time.Now().Add(7 * 24 * time.Hour)
	} else {
		cookieExp = time.Time{}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "userID",
		Value:    signInResponse.User.ID.String(),
		Expires:  cookieExp,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    signInResponse.AccessToken,
		Expires:  cookieExp,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(signInResponse)
}
