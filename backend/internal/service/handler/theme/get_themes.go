package theme

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetThemes(c *fiber.Ctx) error {
	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	// Parse filter parameters
	filter := &models.ThemeFilter{}
	if err := c.QueryParser(filter); err != nil {
		return errs.BadRequest("Invalid Filter Query Parameters")
	}

	// Validate filter parameters
	if validationErrors := xvalidator.Validator.Validate(filter); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	themes, err := h.themeRepository.GetThemes(c.Context(), pagination, filter)
	if err != nil {
		slog.Error("Failed to fetch themes", "error", err)
		return errs.InternalServerError("Failed to fetch themes")
	}

	return c.Status(fiber.StatusOK).JSON(themes)
}
