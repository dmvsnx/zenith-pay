package handlers

import (
	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
	validator *helpers.CustomValidator
}

func NewProductHandler(productUsecase usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		validator:      helpers.NewCustomValidtor(),
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dtos.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.productUsecase.CreateProduct(&req)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Created(c, "Product created successfully", res)
}

func (h *ProductHandler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	res, err := h.productUsecase.GetProductByID(id)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Product retrieved successfully", res)
}

func (h *ProductHandler) ListProduct(c *fiber.Ctx) error {
	p := dtos.PaginationRequest{
		Page:  c.QueryInt("page", dtos.DefaultPage),
		Limit: c.QueryInt("limit", dtos.DefaultLimit),
	}
	p.Normalize()

	res, total, err := h.productUsecase.ListProducts(p.Page, p.Limit)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.PaginatedSuccess(c, "Products retrieved successfully", res, total, p.Page, p.Limit)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dtos.ProductUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return helpers.BadRequest(c, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	res, err := h.productUsecase.UpdateProduct(id, &req)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Product updated successfully", res)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.productUsecase.DeleteProduct(id); err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Product deleted successfully", nil)
}
