package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type CategoryHandler struct {
	categoryUsecase usecase.CategoryUsecase
	validator *helpers.CustomValidator
}

func NewCategoryHandler(categoryUsecase usecase.CategoryUsecase) *CategoryHandler {
	return &CategoryHandler{
		categoryUsecase: categoryUsecase,
		validator: helpers.NewCustomValidtor(),
	}
}

func (h *CategoryHandler) CreateCategory(c *fiber.Ctx) error {
	var req dtos.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.categoryUsecase.CreateCategory(&req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Created(c, "Category created successfully", res)
}

func (h *CategoryHandler) ListCategories(c *fiber.Ctx) error {
	p := dtos.PaginationRequest{
		Page:  c.QueryInt("page", dtos.DefaultPage),
		Limit: c.QueryInt("limit", dtos.DefaultLimit),
	}
	p.Normalize()

	res, total, err := h.categoryUsecase.ListCategories(p.Page, p.Limit)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.PaginatedSuccess(c, "Categories retrieved successfully", res, total, p.Page, p.Limit)
}

func (h *CategoryHandler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")

	res, err := h.categoryUsecase.GetCategoryByID(id)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Category retrieved successfully", res)
}

func (h *CategoryHandler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dtos.CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.categoryUsecase.UpdateCategory(id, &req)
	if err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Category updated successfully", res)
}

func (h *CategoryHandler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.categoryUsecase.DeleteCategory(id); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	return helpers.Success(c, "Category deleted successfully", nil)
}
