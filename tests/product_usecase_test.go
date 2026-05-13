package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/tests/mocks"
)

func TestProductUsecase_CreateProduct_Success(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	now := time.Now()
	productRepo := &mocks.ProductRepo{
		FindByNameFn: func(name string) (*model.Product, error) {
			return nil, nil
		},
		CreateFn: func(product *model.Product) error {
			product.ID = uuid.MustParse("00000000-0000-0000-0000-000000000011")
			product.CreatedAt = now
			product.UpdatedAt = now
			return nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{ID: catID, Name: "Food"}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	res, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: catID.String(),
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Name != "Burger" {
		t.Fatalf("expected Burger, got %s", res.Name)
	}
	if res.CategoryName != "Food" {
		t.Fatalf("expected Food, got %s", res.CategoryName)
	}
}

func TestProductUsecase_CreateProduct_InvalidCategoryID(t *testing.T) {
	uc := usecase.NewProductUsecase(&mocks.ProductRepo{}, &mocks.CategoryRepo{})

	_, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: "not-a-uuid",
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err == nil || err.Error() != "invalid category ID" {
		t.Fatalf("expected invalid category ID, got %v", err)
	}
}

func TestProductUsecase_CreateProduct_NilCategoryID(t *testing.T) {
	uc := usecase.NewProductUsecase(&mocks.ProductRepo{}, &mocks.CategoryRepo{})

	_, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: uuid.Nil.String(),
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err == nil || err.Error() != "category ID cannot be empty" {
		t.Fatalf("expected category ID cannot be empty, got %v", err)
	}
}

func TestProductUsecase_CreateProduct_NameExists(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByNameFn: func(name string) (*model.Product, error) {
			return &model.Product{Name: name}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	_, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: "00000000-0000-0000-0000-000000000010",
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err == nil || err.Error() != "product with the same name already exists" {
		t.Fatalf("expected product with the same name already exists, got %v", err)
	}
}

func TestProductUsecase_CreateProduct_CategoryNotFound(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByNameFn: func(name string) (*model.Product, error) {
			return nil, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	_, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: catID.String(),
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestProductUsecase_CreateProduct_FailCreate(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByNameFn: func(name string) (*model.Product, error) {
			return nil, nil
		},
		CreateFn: func(product *model.Product) error {
			return errors.New("db error")
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{ID: catID, Name: "Food"}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	_, err := uc.CreateProduct(&dtos.ProductRequest{
		CategoryID: catID.String(),
		Name:       "Burger",
		Price:      25000,
		Stock:      10,
	})

	if err == nil || err.Error() != "failed to create product" {
		t.Fatalf("expected failed to create product, got %v", err)
	}
}

func TestProductUsecase_GetProductByID_Success(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	now := time.Now()
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
				Name:       "Burger",
				Price:      25000,
				Stock:      10,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{ID: catID, Name: "Food"}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	res, err := uc.GetProductByID("00000000-0000-0000-0000-000000000011")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Name != "Burger" {
		t.Fatalf("expected Burger, got %s", res.Name)
	}
}

func TestProductUsecase_GetProductByID_NotFound(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	_, err := uc.GetProductByID("00000000-0000-0000-0000-000000000011")

	if err == nil || err.Error() != "product not found" {
		t.Fatalf("expected product not found, got %v", err)
	}
}

func TestProductUsecase_GetProductByID_CategoryNotFound(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
			}, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	_, err := uc.GetProductByID("00000000-0000-0000-0000-000000000011")

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestProductUsecase_ListProducts_Success(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	now := time.Now()
	productRepo := &mocks.ProductRepo{
		FindAllPaginatedFn: func(offset, limit int) ([]*model.Product, int64, error) {
			return []*model.Product{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000011"), CategoryID: catID, Name: "Burger", Price: 25000, Stock: 10, CreatedAt: now, UpdatedAt: now},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000012"), CategoryID: catID, Name: "Fries", Price: 15000, Stock: 20, CreatedAt: now, UpdatedAt: now},
			}, 2, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{ID: catID, Name: "Food"}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	res, total, err := uc.ListProducts(1, 10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 products, got %d", len(res))
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
}

func TestProductUsecase_ListProducts_FailFetch(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindAllPaginatedFn: func(offset, limit int) ([]*model.Product, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	_, _, err := uc.ListProducts(1, 10)

	if err == nil || err.Error() != "failed to retrieve products" {
		t.Fatalf("expected failed to retrieve products, got %v", err)
	}
}

func TestProductUsecase_ListProducts_CategoryNotFound(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	now := time.Now()
	productRepo := &mocks.ProductRepo{
		FindAllPaginatedFn: func(offset, limit int) ([]*model.Product, int64, error) {
			return []*model.Product{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000011"), CategoryID: catID, Name: "Burger", Price: 25000, Stock: 10, CreatedAt: now, UpdatedAt: now},
			}, 1, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	_, _, err := uc.ListProducts(1, 10)

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestProductUsecase_UpdateProduct_Success(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
				Name:       "Burger",
				Price:      25000,
				Stock:      10,
			}, nil
		},
		UpdateFn: func(product *model.Product) error {
			return nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	newName := "Cheese Burger"
	newPrice := int64(30000)
	err := uc.UpdateProduct("00000000-0000-0000-0000-000000000011", &dtos.ProductUpdateRequest{
		Name:  &newName,
		Price: &newPrice,
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductUsecase_UpdateProduct_NotFound(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	newName := "Cheese Burger"
	err := uc.UpdateProduct("00000000-0000-0000-0000-000000000011", &dtos.ProductUpdateRequest{
		Name: &newName,
	})

	if err == nil || err.Error() != "product not found" {
		t.Fatalf("expected product not found, got %v", err)
	}
}

func TestProductUsecase_UpdateProduct_InvalidCategoryID(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
			}, nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	badID := "not-a-uuid"
	err := uc.UpdateProduct("00000000-0000-0000-0000-000000000011", &dtos.ProductUpdateRequest{
		CategoryID: &badID,
	})

	if err == nil || err.Error() != "invalid category ID" {
		t.Fatalf("expected invalid category ID, got %v", err)
	}
}

func TestProductUsecase_UpdateProduct_CategoryNotFound(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
			}, nil
		},
	}
	categoryRepo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, categoryRepo)

	newCatID := "00000000-0000-0000-0000-000000000099"
	err := uc.UpdateProduct("00000000-0000-0000-0000-000000000011", &dtos.ProductUpdateRequest{
		CategoryID: &newCatID,
	})

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestProductUsecase_UpdateProduct_FailUpdate(t *testing.T) {
	catID := uuid.MustParse("00000000-0000-0000-0000-000000000010")
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID:         uuid.MustParse("00000000-0000-0000-0000-000000000011"),
				CategoryID: catID,
				Name:       "Burger",
			}, nil
		},
		UpdateFn: func(product *model.Product) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	newName := "Cheese Burger"
	err := uc.UpdateProduct("00000000-0000-0000-0000-000000000011", &dtos.ProductUpdateRequest{
		Name: &newName,
	})

	if err == nil || err.Error() != "failed to update product" {
		t.Fatalf("expected failed to update product, got %v", err)
	}
}

func TestProductUsecase_DeleteProduct_Success(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000011"),
			}, nil
		},
		DeleteFn: func(id string) error {
			return nil
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	err := uc.DeleteProduct("00000000-0000-0000-0000-000000000011")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestProductUsecase_DeleteProduct_NotFound(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	err := uc.DeleteProduct("00000000-0000-0000-0000-000000000011")

	if err == nil || err.Error() != "product not found" {
		t.Fatalf("expected product not found, got %v", err)
	}
}

func TestProductUsecase_DeleteProduct_FailDelete(t *testing.T) {
	productRepo := &mocks.ProductRepo{
		FindByIDFn: func(id string) (*model.Product, error) {
			return &model.Product{
				ID: uuid.MustParse("00000000-0000-0000-0000-000000000011"),
			}, nil
		},
		DeleteFn: func(id string) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewProductUsecase(productRepo, &mocks.CategoryRepo{})

	err := uc.DeleteProduct("00000000-0000-0000-0000-000000000011")

	if err == nil || err.Error() != "failed to delete product" {
		t.Fatalf("expected failed to delete product, got %v", err)
	}
}
