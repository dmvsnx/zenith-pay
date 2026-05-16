package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
	validator *helpers.CustomValidator
}

func NewTransactionHandler(tu usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: tu,
		validator: helpers.NewCustomValidtor(),
	}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	var req dtos.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	userID, ok := c.Locals("userID").(string)
	if !ok || userID == "" {
		return helpers.Unauthorized(c, "Unauthorized")
	}

	shiftID, ok := c.Locals("shiftID").(string)
	if !ok || shiftID == "" {
		return helpers.Forbidden(c, "No active shift")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.transactionUsecase.CreateTransaction(userID, shiftID, &req)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Created(c, "Transaction created successfully", res)
}

func (h *TransactionHandler) ListTransactions(c *fiber.Ctx) error {
	p := dtos.PaginationRequest{
		Page:  c.QueryInt("page", dtos.DefaultPage),
		Limit: c.QueryInt("limit", dtos.DefaultLimit),
	}
	p.Normalize()

	res, total, err := h.transactionUsecase.GetAllTransaction(p.Page, p.Limit)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.PaginatedSuccess(c, "Transactions retrieved successfully", res, total, p.Page, p.Limit)
}

func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.transactionUsecase.GetTransactionByID(id)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Transaction retrieved successfully", res)
}
