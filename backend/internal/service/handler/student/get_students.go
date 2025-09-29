package student

import (
	"specialstandard/internal/errs"
	"specialstandard/internal/utils"
	"specialstandard/internal/xvalidator"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetStudents(c *fiber.Ctx) error {
	pagination := utils.NewPagination()
	if err := c.QueryParser(&pagination); err != nil {
		return errs.BadRequest("Invalid Pagination Query Parameters")
	}

	if validationErrors := xvalidator.Validator.Validate(pagination); len(validationErrors) > 0 {
		return errs.InvalidRequestData(xvalidator.ConvertToMessages(validationErrors))
	}

	students, err := h.studentRepository.GetStudents(c.Context(), pagination)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(students)
}
