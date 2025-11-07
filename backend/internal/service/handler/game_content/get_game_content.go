package game_content

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetGameContents(c *fiber.Ctx) error {
	var getGameContentReq models.GetGameContentRequest
	if err := c.QueryParser(&getGameContentReq); err != nil {
		return errs.BadRequest("GameContent Query-Parameters Parsing Error")
	}

	if validationErrors := xvalidator.Validator.Validate(getGameContentReq); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	gameContents, err := h.gameContentRepository.GetGameContents(c.Context(), getGameContentReq)
	if err != nil {
		req := getGameContentReq
		// For all other database errors, return internal server error without exposing details
		slog.Error("Failed to get game contents", "theme_id", req.ThemeID, "category",
			req.Category, "question_type", req.QuestionType, "difficulty_level",
			req.DifficultyLevel, "count", req.Count)
		return errs.InternalServerError("Failed to retrieve game contents")
	}

	return c.Status(fiber.StatusOK).JSON(gameContents)
}
