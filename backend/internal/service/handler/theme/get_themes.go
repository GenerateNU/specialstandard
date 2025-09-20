package theme

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetThemes(c *fiber.Ctx) error {
	themes, err := h.themeRepository.GetThemes(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(themes)
}
