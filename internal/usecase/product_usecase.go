package usecase

import (
	"errors"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/utils"
)

type ProductUsecase interface {
	CreateProduct(req *dtos.ProductRequest) (*dtos.ProductResponse, error)
	GetProductByID(id string) (*dtos.ProductResponse, error)
	ListProducts(page, limit int) ([]*dtos.ProductResponse, int64, error)
	UpdateProduct(id string, req *dtos.ProductUpdateRequest) error
	DeleteProduct(id string) error
}

type productUsecase struct {
	productRepo repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductUsecase(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductUsecase {
	return &productUsecase{
		productRepo: productRepo,
		categoryRepo: categoryRepo,
	}
}

func (u *productUsecase) CreateProduct(req *dtos.ProductRequest) (*dtos.ProductResponse, error) {
	categoryID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return nil, errors.New("invalid category ID")
	}
	if categoryID == uuid.Nil {
		return nil, errors.New("category ID cannot be empty")
	}

	existingProduct, err := u.productRepo.FindByName(req.Name)
	if err == nil && existingProduct != nil {
		return nil, errors.New("product with the same name already exists")
	}

	category, err := u.categoryRepo.FindByID(req.CategoryID)
	if err != nil || category == nil {
		return nil, errors.New("category not found")
	}

	sku, err := utils.GenerateSKU()
	if err != nil {
		return nil, errors.New("failed to generate SKU")
	}

	product := &model.Product{
		CategoryID: categoryID,
		SKU:        sku,
		Name:       req.Name,
		Price:      req.Price,
		Stock:      req.Stock,
	}

	if err := u.productRepo.Create(product); err != nil {
		return nil, errors.New("failed to create product")
	}

	res := &dtos.ProductResponse{
		ID:         product.ID.String(),
		CategoryID: product.CategoryID.String(),
		CategoryName: category.Name,
		SKU:        product.SKU,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CreatedAt:  product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return res, nil
}

func (u *productUsecase) GetProductByID(id string) (*dtos.ProductResponse, error) {
	product, err := u.productRepo.FindByID(id)
	if err != nil || product == nil {
		return nil, errors.New("product not found")
	}

	category, err := u.categoryRepo.FindByID(product.CategoryID.String())
	if err != nil || category == nil {
		return nil, errors.New("category not found")
	}

	res := &dtos.ProductResponse{
		ID:         product.ID.String(),
		CategoryID: product.CategoryID.String(),
		CategoryName: category.Name,
		SKU:        product.SKU,
		Name:       product.Name,
		Price:      product.Price,
		Stock:      product.Stock,
		CreatedAt:  product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return res, nil
}

func (u *productUsecase) ListProducts(page, limit int) ([]*dtos.ProductResponse, int64, error) {
	products, total, err := u.productRepo.FindAllPaginated((page-1)*limit, limit)
	if err != nil {
		return nil, 0, errors.New("failed to retrieve products")
	}

	res := make([]*dtos.ProductResponse, 0, len(products))

	for _, product := range products {
		category, err := u.categoryRepo.FindByID(product.CategoryID.String())
		if err != nil || category == nil {
			return nil, 0, errors.New("category not found")
		}

		res = append(res, &dtos.ProductResponse{
			ID:         product.ID.String(),
			CategoryID: product.CategoryID.String(),
			CategoryName: category.Name,
			SKU:        product.SKU,
			Name:       product.Name,
			Price:      product.Price,
			Stock:      product.Stock,
			CreatedAt:  product.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  product.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return res, total, nil
}

func (u *productUsecase) UpdateProduct(id string, req *dtos.ProductUpdateRequest) error {
	product, err := u.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if req.CategoryID != nil {
		if _, err := uuid.Parse(*req.CategoryID); err != nil {
			return errors.New("invalid category ID")
		}

		category, err := u.categoryRepo.FindByID(*req.CategoryID)
		if err != nil {
			return errors.New("category not found")
		}
		product.CategoryID = category.ID
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	if err := u.productRepo.Update(product); err != nil {
		return errors.New("failed to update product")
	}

	return nil
}

func (u *productUsecase) DeleteProduct(id string) error {
	_, err := u.productRepo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	if err := u.productRepo.Delete(id); err != nil {
		return errors.New("failed to delete product")
	}

	return nil
}
