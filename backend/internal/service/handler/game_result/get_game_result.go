package game_result

import (
	"log/slog"
	"specialstandard/internal/errs"
	"specialstandard/internal/models"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetGameResults(c *fiber.Ctx) error {
	var gameResultsReq models.GetGameResultQuery
	if err := c.QueryParser(&gameResultsReq); err != nil {
		return errs.BadRequest("GameResults Query-Parameters Parsing Error")
	}

	if validationErrors := xvalidator.Validator.Validate(gameResultsReq); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	gameResults, err := h.gameResultRepository.GetGameResults(c.Context(), &gameResultsReq, pagination)
	if err != nil {
		// For all database errors, return internal server error without exposing details
		slog.Error("Failed to get session", "err", err)
		return errs.InternalServerError("Failed to retrieve sessions")
	}

	return c.Status(fiber.StatusOK).JSON(gameResults)
}
