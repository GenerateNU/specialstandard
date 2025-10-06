package auth

import (
	"fmt"
	"net/http"
	"specialstandard/internal/config"

	"github.com/gofiber/fiber/v2"
)

func Middleware(cfg *config.Supabase) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("jwt", "")

		if token == "" {
			fmt.Println("JWT Not Found in Middleware")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token Not Found"})
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth/v1/user", cfg.URL), nil)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("apikey", cfg.ServiceRoleKey)

		res, err := Client.Do(req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to validate token"})
		}
		defer func() {
			_ = res.Body.Close()
		}()

		if res.StatusCode != http.StatusOK {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid/Expired Token"})
		}

		return c.Next()
	}
}
