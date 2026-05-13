package tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/tests/mocks"
)

func TestCategoryUsecase_CreateCategory_Success(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByNameFn: func(name string) (*model.Category, error) {
			return nil, nil
		},
		CreateFn: func(category *model.Category) error {
			category.ID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
			return nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	res, err := uc.CreateCategory(&dtos.CategoryRequest{Name: "Food"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Name != "Food" {
		t.Fatalf("expected Food, got %s", res.Name)
	}
}

func TestCategoryUsecase_CreateCategory_AlreadyExists(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByNameFn: func(name string) (*model.Category, error) {
			return &model.Category{Name: name}, nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.CreateCategory(&dtos.CategoryRequest{Name: "Food"})

	if err == nil || err.Error() != "category already exists" {
		t.Fatalf("expected category already exists, got %v", err)
	}
}

func TestCategoryUsecase_CreateCategory_FailCreate(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByNameFn: func(name string) (*model.Category, error) {
			return nil, nil
		},
		CreateFn: func(category *model.Category) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.CreateCategory(&dtos.CategoryRequest{Name: "Food"})

	if err == nil || err.Error() != "failed to create category" {
		t.Fatalf("expected failed to create category, got %v", err)
	}
}

func TestCategoryUsecase_ListCategories_Success(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindAllPaginatedFn: func(offset, limit int) ([]*model.Category, int64, error) {
			return []*model.Category{
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"), Name: "Food"},
				{ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"), Name: "Drink"},
			}, 2, nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	res, total, err := uc.ListCategories(1, 10)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(res))
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
}

func TestCategoryUsecase_ListCategories_FailFetch(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindAllPaginatedFn: func(offset, limit int) ([]*model.Category, int64, error) {
			return nil, 0, errors.New("db error")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, _, err := uc.ListCategories(1, 10)

	if err == nil || err.Error() != "failed to fetch categories" {
		t.Fatalf("expected failed to fetch categories, got %v", err)
	}
}

func TestCategoryUsecase_GetCategoryByID_Success(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name: "Food",
			}, nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	res, err := uc.GetCategoryByID("00000000-0000-0000-0000-000000000001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Name != "Food" {
		t.Fatalf("expected Food, got %s", res.Name)
	}
}

func TestCategoryUsecase_GetCategoryByID_NotFound(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.GetCategoryByID("00000000-0000-0000-0000-000000000001")

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestCategoryUsecase_UpdateCategory_Success(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name: "Old Name",
			}, nil
		},
		UpdateFn: func(category *model.Category) error {
			return nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	res, err := uc.UpdateCategory("00000000-0000-0000-0000-000000000001", &dtos.CategoryRequest{Name: "New Name"})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Name != "New Name" {
		t.Fatalf("expected New Name, got %s", res.Name)
	}
}

func TestCategoryUsecase_UpdateCategory_NotFound(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.UpdateCategory("00000000-0000-0000-0000-000000000001", &dtos.CategoryRequest{Name: "New Name"})

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestCategoryUsecase_UpdateCategory_FailUpdate(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name: "Old Name",
			}, nil
		},
		UpdateFn: func(category *model.Category) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	_, err := uc.UpdateCategory("00000000-0000-0000-0000-000000000001", &dtos.CategoryRequest{Name: "New Name"})

	if err == nil || err.Error() != "failed to update category" {
		t.Fatalf("expected failed to update category, got %v", err)
	}
}

func TestCategoryUsecase_DeleteCategory_Success(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name: "Food",
			}, nil
		},
		DeleteFn: func(id string) error {
			return nil
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	err := uc.DeleteCategory("00000000-0000-0000-0000-000000000001")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCategoryUsecase_DeleteCategory_NotFound(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return nil, errors.New("not found")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	err := uc.DeleteCategory("00000000-0000-0000-0000-000000000001")

	if err == nil || err.Error() != "category not found" {
		t.Fatalf("expected category not found, got %v", err)
	}
}

func TestCategoryUsecase_DeleteCategory_FailDelete(t *testing.T) {
	repo := &mocks.CategoryRepo{
		FindByIDFn: func(id string) (*model.Category, error) {
			return &model.Category{
				ID:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
				Name: "Food",
			}, nil
		},
		DeleteFn: func(id string) error {
			return errors.New("db error")
		},
	}
	uc := usecase.NewCategoryUsecase(repo)

	err := uc.DeleteCategory("00000000-0000-0000-0000-000000000001")

	if err == nil || err.Error() != "failed to delete category" {
		t.Fatalf("expected failed to delete category, got %v", err)
	}
}
