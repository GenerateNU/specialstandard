package session

import (
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetSessions(c *fiber.Ctx) error {
	sessions, err := h.sessionRepository.GetSessions(c.Context())
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session", "err", err)
		return errs.InternalServerError("Failed to retrieve sessions")
	}

	return c.Status(fiber.StatusOK).JSON(sessions)
}
