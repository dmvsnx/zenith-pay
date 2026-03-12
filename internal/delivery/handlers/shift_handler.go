package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type ShiftHandler struct {
	usecase usecase.ShiftUsecase
	validator *helpers.CustomValidator
}

func NewShiftHandler(uc usecase.ShiftUsecase) *ShiftHandler {
	return &ShiftHandler{
		usecase: uc,
	}
}

func (h *ShiftHandler) OpenShift(c *fiber.Ctx) error {
	cashierID := c.Locals("userID").(string)
	var req dtos.OpenShiftRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.usecase.OpenShift(cashierID, req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Shift opened successfully", res)
}

func (h *ShiftHandler) CloseShift(c *fiber.Ctx) error {
	cashierID := c.Locals("userID").(string)
	var req dtos.CloseShiftRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.usecase.CloseShift(cashierID, req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Shift closed successfully", res)
}

func (h *ShiftHandler) GetActiveShift(c *fiber.Ctx) error {
	cashierID := c.Locals("userID").(string)

	res, err := h.usecase.GetActiveShift(cashierID)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Shift retrieved successfully", res)
}
