package game_content

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (h *Handler) GetGameContents(c *fiber.Ctx) error {
	var getGameContentReq models.GetGameContentRequest
	if err := c.QueryParser(&getGameContentReq); err != nil {
		return errs.BadRequest("GameContent Query-Parameters Parsing Error")
	}

	if validationErrors := xvalidator.Validator.Validate(getGameContentReq); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	gameContent, err := h.gameContentRepository.GetGameContents(c.Context(), getGameContentReq)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("Game Contents not found")
		}

		// For all other database errors, return internal server error without exposing details
		slog.Error("Failed to get game contents", "category", getGameContentReq.Category, "level",
			getGameContentReq.Level, "count", getGameContentReq.Count)
		return errs.InternalServerError("Failed to retrieve session")
	}

	return c.Status(fiber.StatusOK).JSON(gameContent)
}
