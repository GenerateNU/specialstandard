package theme

import (
	"specialstandard/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CreateTheme(c *fiber.Ctx) error {
	var theme models.CreateThemeInput
	if err := c.BodyParser(&theme); err != nil {
		return err
	}

	createdTheme, err := h.themeRepository.CreateTheme(c.Context(), &theme)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(createdTheme)
}
