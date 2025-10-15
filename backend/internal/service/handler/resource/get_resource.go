package resource

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"
	"strconv"
	"strings"
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
	themeName := c.Query("theme_name")

	themeMonthStr := c.Query("theme_month")
	themeYearStr := c.Query("theme_year")

	var themeMonth, themeYear *int
	if themeMonthStr != "" {
		parsedMonth, err := strconv.Atoi(themeMonthStr)
		if err == nil {
			if parsedMonth < 1 || parsedMonth > 12 {
				return errs.InvalidRequestData(map[string]string{"theme_month": "month must be between 1 and 12"})
			}
			themeMonth = &parsedMonth
		}
	}
	if themeYearStr != "" {
		parsedYear, err := strconv.Atoi(themeYearStr)
		if err == nil {
			if parsedYear < 1900 || parsedYear > 3000 {
				return errs.InvalidRequestData(map[string]string{"theme_year": "year is invalid"})
			}
			themeYear = &parsedYear
		}
	}

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	resourcesWithThemes, err := h.resourceRepository.GetResources(c.Context(), themeId, gradeLevel, res_type, title, category, content, themeName, date, themeMonth, themeYear, pagination)
	if err != nil {
		return errs.InternalServerError(err.Error())
	}

	var resources []models.ResourceResponseWithURL
	for _, res := range resourcesWithThemes {
		key := strings.TrimPrefix(*res.Content, "/")
		presignedURL := ""

		if key != "" {
			url, err := h.s3Client.GeneratePresignedURL(c.Context(), key, 15*time.Minute)
			if err != nil {
				slog.Warn("Failed to generate presigned URL for resource",
					"key", key,
					"error", err,
				)
			} else {
				presignedURL = url
			}
		}
		resources = append(resources, models.ResourceResponseWithURL{
			ResourceWithTheme: res,
			PresignedURL:      presignedURL,
		})
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
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			return errs.NotFound("resource", "resource not found")
		}
		return errs.InternalServerError()
	}

	key := strings.TrimPrefix(*resource.Content, "/")
	presignedURL := ""

	if key != "" {
		url, err := h.s3Client.GeneratePresignedURL(c.Context(), key, 15*time.Minute)
		if err != nil {
			slog.Warn("Failed to generate presigned URL for resource,",
				"key", key,
				"error", err)
		} else {
			presignedURL = url
		}
	}

	response := models.ResourceResponseWithURL{
		ResourceWithTheme: *resource,
		PresignedURL:      presignedURL,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
