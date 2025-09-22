package session

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	sessions, err := h.sessionRepository.GetSessions(c.Context(), pagination)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session", "err", err)
		return errs.InternalServerError("Failed to retrieve sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
