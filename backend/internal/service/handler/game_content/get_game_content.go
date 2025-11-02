package game_content

import (
	"errors"
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

func (h *Handler) GetGameContents(c *fiber.Ctx) error {
	var getGameContentReq models.GetGameContentRequest
	if err := c.QueryParser(&getGameContentReq); err != nil {
		return errs.BadRequest("GameContent Query-Parameters Parsing Error")
	}

	if validationErrors := xvalidator.Validator.Validate(getGameContentReq); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	gameContent, err := h.gameContentRepository.GetGameContent(c.Context(), getGameContentReq)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errs.NotFound("Game Contents Not Found")
		}

		// For all other database errors, return internal server error without exposing details
		slog.Error("Failed to get game contents", "category", getGameContentReq.Category, "level",
			getGameContentReq.Level, "count", getGameContentReq.Count)
		return errs.InternalServerError("Failed to retrieve game contents")
	}

	return c.Status(fiber.StatusOK).JSON(gameContent)
}
