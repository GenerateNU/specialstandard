package game_result

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/xvalidator"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) PostGameResults(c *fiber.Ctx) error {
	var postGameResult models.PostGameResult

	if err := c.BodyParser(&postGameResult); err != nil {
		return errs.InvalidJSON("Failed to parse PostGameResult data")
	}

	if validationErrors := h.validator.Validate(postGameResult); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	newGameResult, err := h.gameResultRepository.PostGameResult(c.Context(), postGameResult)
	if err != nil {
		slog.Error("Failed to post game-result", "err", err)
		errStr := err.Error()
		switch {
		case strings.Contains(errStr, "foreign key"):
			return errs.BadRequest("Invalid Reference")
		case strings.Contains(errStr, "check constraint"):
			return errs.BadRequest("Violated a check constraint")
		case strings.Contains(errStr, "connection refused"):
			return errs.InternalServerError("Database Connection Error")
		default:
			return errs.InternalServerError("Failed to Create GameResult")
		}
	}

	return c.Status(fiber.StatusCreated).JSON(newGameResult)
}
