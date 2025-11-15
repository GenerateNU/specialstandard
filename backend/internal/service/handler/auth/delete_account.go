package auth

import (
	"fmt"
	"specialstandard/internal/auth"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

// DeleteAccount handler to revoke session, delete the account and clear cookies
func (h *Handler) DeleteAccount(c *fiber.Ctx, id string) error {
	// Retrieve the JWT token from cookies
	accessToken := c.Cookies("jwt")
	if accessToken == "" {
		return errs.Unauthorized("No authentication token found")
	}

	// Verify that the user ID from the token matches the requested ID to delete
	// This is a security check to prevent users from deleting other accounts
	userID := c.Cookies("userID")
	if userID != id {
		return errs.Forbidden("You can only delete your own account")
	}

	err := h.therapistRepository.DeleteTherapist(c.Context(), id)
	if err != nil {
		fmt.Println("Error deleting user from database:", err)
		return errs.InternalServerError(fmt.Sprintf("Failed to delete user data: %v", err))
	}

	err = auth.SupabaseDeleteAccount(&h.config, userID)
	if err != nil {
		fmt.Println("Error deleting account from Supabase:", err)
		return errs.InternalServerError(fmt.Sprintf("Failed to delete account: %v", err))
	}

	// Clear cookies related to authentication
	c.ClearCookie("jwt")
	c.ClearCookie("userID")
	c.ClearCookie("refreshToken")

	// Return success response
	return c.SendStatus(fiber.StatusOK)
}
