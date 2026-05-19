package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
	validator *helpers.CustomValidator
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
		validator: helpers.NewCustomValidtor(),
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dtos.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.userUsecase.Register(&req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Created(c, "User registered successfully", res)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	p := dtos.PaginationRequest{
		Page:  c.QueryInt("page", dtos.DefaultPage),
		Limit: c.QueryInt("limit", dtos.DefaultLimit),
	}
	p.Normalize()

	res, total, err := h.userUsecase.ListUsers(p.Page, p.Limit)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.PaginatedSuccess(c, "Users retrieved successfully", res, total, p.Page, p.Limit)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dtos.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.userUsecase.Login(&req)
	if err != nil {
		if err.Error() == "invalid username or password" {
			return helpers.Unauthorized(c, err.Error())
		}
		return helpers.InternalServerError(c, "internal server error")
	}

	return helpers.Success(c, "Login successful", res)
}
