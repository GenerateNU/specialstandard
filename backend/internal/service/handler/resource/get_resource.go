package resource

import (
	"specialstandard/internal/errs"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetResources(c *fiber.Ctx) error {
	themeIdStr := c.Query("theme_id")
	var themeId uuid.UUID
	if themeIdStr != "" {
		parsedThemeId, err := uuid.Parse(themeIdStr)
		if err != nil {
			return errs.InvalidRequestData(map[string]string{"theme_id": "invalid UUID"})
		}
		themeId = parsedThemeId
	}
	gradeLevel := c.Query("grade_level")
	dateStr := c.Query("date")
	var date *time.Time
	if dateStr != "" {
		parsedDate, err := time.Parse("2006-01-02", dateStr)
		if err == nil {
			date = &parsedDate
		}
	}
	res_type := c.Query("type")
	title := c.Query("title")
	category := c.Query("category")
	content := c.Query("content")

	resources, err := h.resourceRepository.GetResources(c.Context(), themeId, gradeLevel, res_type, title, category, content, date)
	if err != nil {
		return errs.InternalServerError(err.Error())
	}

	return c.JSON(resources)
}

func (h *Handler) GetResourceByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.InvalidRequestData(map[string]string{"id": "invalid resource UUID"})
	}

	resource, err := h.resourceRepository.GetResourceByID(c.Context(), id)
	if resource == nil {
		return errs.NotFound("resource", "resource not found")
	}
	if err != nil {
		return errs.InternalServerError()
	}
	return c.Status(fiber.StatusOK).JSON(resource)
}
