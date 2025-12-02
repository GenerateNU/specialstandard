package session

import (
	"fmt"
	"log/slog"
	"specialstandard/internal/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) DeleteRecurringSessions(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return errs.BadRequest("Parsing Error: Invalid ID Format. ID: " + id.String())
	}

	err = h.sessionRepository.DeleteRecurringSessions(c.Context(), id)
	if err != nil {
		slog.Error("Failed to delete sessions", "id", id, "err", err)
		return errs.InternalServerError("Internal Server Error")
	}

	fmt.Println("Deleted recurring sessions with ID:", id)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Sessions deleted successfully",
	})
}
