package game_content

import (
	"context"
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetGameContents(c *fiber.Ctx) error {
	getGameContentReq := models.NewGetGameContentRequest()
	if err := c.QueryParser(&getGameContentReq); err != nil {
		return errs.BadRequest("GameContent Query-Parameters Parsing Error")
	}

	if validationErrors := xvalidator.Validator.Validate(getGameContentReq); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	gameContents, err := h.gameContentRepository.GetGameContents(c.Context(), getGameContentReq)
	if err != nil {
		req := getGameContentReq
		slog.Error("Failed to get game contents", "theme_id", req.ThemeID, "category",
			req.Category, "question_type", req.QuestionType, "difficulty_level",
			req.DifficultyLevel, "question_count", req.QuestionCount, "words_count",
			req.WordsCount)
		return errs.InternalServerError("Failed to retrieve game contents")
	}

	// Generate presigned URLs for answer field
	if h.s3Client != nil {
		for i := range gameContents {
			if gameContents[i].Answer != "" {
				presignedURL, err := h.s3Client.GeneratePresignedURL(context.Background(), gameContents[i].Answer, time.Hour)
				if err != nil {
					slog.Warn("Failed to generate presigned URL", "key", gameContents[i].Answer, "error", err)
				} else {
					gameContents[i].Answer = presignedURL
				}
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(gameContents)
}
