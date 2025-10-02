package auth

import (
	"fmt"
	"log/slog"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func derefString(s *string) string {
	return *s
}

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var cred Credentials

	if err := c.BodyParser(&cred); err != nil {
		return errs.BadRequest(fmt.Sprintf("Invalid Request Body: %v", cred))
	}

	res, err := auth.SupabaseSignup(&h.config, cred.Email, cred.Password)
	if err != nil {
		slog.Error(fmt.Sprintf("Signup Request Failed: %v", err))
		return errs.InternalServerError(fmt.Sprintf("Signup Request Failed: %v", err))
	}

	postTherapist := models.CreateTherapistInput{
		ID:        res.User.ID,
		FirstName: derefString(cred.FirstName),
		LastName:  derefString(cred.LastName),
		Email:     cred.Email,
	}
	_, err = h.therapistRepository.CreateTherapist(c.Context(), &postTherapist)
	if err != nil {
		return errs.BadRequest(fmt.Sprintf("Creating Therapist/User failed: %v", err))
	}

	expiration := time.Now().Add(30 * 24 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    res.AccessToken,
		Expires:  expiration,
		Secure:   true,
		SameSite: "Lax",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "userID",
		Value:    res.User.ID.String(),
		Expires:  expiration,
		Secure:   true,
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusCreated).JSON(res)
}
