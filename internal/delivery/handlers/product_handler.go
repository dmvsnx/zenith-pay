package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/storage/minio"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type ProductHandler struct {
	productUsecase usecase.ProductUsecase
	minioService   minio.Service
	validator      *helpers.CustomValidator
}

func NewProductHandler(productUsecase usecase.ProductUsecase, minioService minio.Service) *ProductHandler {
	return &ProductHandler{
		productUsecase: productUsecase,
		minioService:   minioService,
		validator:      helpers.NewCustomValidtor(),
	}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return helpers.BadRequest(c, "Invalid multipart form data")
	}

	price, err := strconv.ParseInt(c.FormValue("price"), 10, 64)
	if err != nil {
		return helpers.BadRequest(c, "Invalid price")
	}

	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return helpers.BadRequest(c, "Invalid stock")
	}

	req := dtos.ProductRequest{
		CategoryID: c.FormValue("category_id"),
		Name:       c.FormValue("name"),
		Price:      price,
		Stock:      stock,
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	fileHeaders := form.File["image"]
	if len(fileHeaders) == 0 {
		return helpers.BadRequest(c, "Image is required")
	}

	imageURL, err := h.minioService.UploadImage(fileHeaders[0])
	if err != nil {
		return helpers.InternalServerError(c, "Failed to upload image")
	}
	req.Image = imageURL

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

	form, err := c.MultipartForm()
	if err != nil {
		return helpers.BadRequest(c, "Invalid multipart form data")
	}

	req := dtos.ProductUpdateRequest{}

	if v := c.FormValue("category_id"); v != "" {
		req.CategoryID = &v
	}
	if v := c.FormValue("name"); v != "" {
		req.Name = &v
	}
	if v := c.FormValue("price"); v != "" {
		price, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return helpers.BadRequest(c, "Invalid price")
		}
		req.Price = &price
	}
	if v := c.FormValue("stock"); v != "" {
		stock, err := strconv.Atoi(v)
		if err != nil {
			return helpers.BadRequest(c, "Invalid stock")
		}
		req.Stock = &stock
	}

	if err := h.validator.Validate(&req); err != nil {
		return helpers.BadRequest(c, err.Error())
	}

	fileHeaders := form.File["image"]
	if len(fileHeaders) > 0 {
		imageURL, err := h.minioService.UploadImage(fileHeaders[0])
		if err != nil {
			return helpers.InternalServerError(c, "Failed to upload image")
		}
		req.Image = &imageURL
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
