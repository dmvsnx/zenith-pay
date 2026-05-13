package usecase

import (
	"errors"

	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
)

type CategoryUsecase interface {
	CreateCategory(req *dtos.CategoryRequest) (*dtos.CategoryResponse, error)
	ListCategories(page, limit int) ([]*dtos.CategoryResponse, int64, error)
	GetCategoryByID(id string) (*dtos.CategoryResponse, error)
	UpdateCategory(id string, req *dtos.CategoryRequest) (*dtos.CategoryResponse, error)
	DeleteCategory(id string) error
}

type categoryUsecase struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryUsecase(categoryRepo repository.CategoryRepository) CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: categoryRepo,
	}
}

func (u *categoryUsecase) CreateCategory(req *dtos.CategoryRequest) (*dtos.CategoryResponse, error) {
	existstingCategory, err := u.categoryRepo.FindByName(req.Name)
	if err == nil && existstingCategory != nil {
		return nil, errors.New("category already exists")
	}

	category := &model.Category{
		Name: req.Name,
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, errors.New("failed to create category")
	}

	res := &dtos.CategoryResponse{
		ID:  category.ID.String(),
		Name: category.Name,
	}

	return res, nil
}

func (u *categoryUsecase) ListCategories(page, limit int) ([]*dtos.CategoryResponse, int64, error) {
	categories, total, err := u.categoryRepo.FindAllPaginated((page-1)*limit, limit)
	if err != nil {
		return nil, 0, errors.New("failed to fetch categories")
	}

	var res []*dtos.CategoryResponse
	for _, category := range categories {
		res = append(res, &dtos.CategoryResponse{
			ID:  category.ID.String(),
			Name: category.Name,
		})
	}

	return res, total, nil
}

func (u *categoryUsecase) GetCategoryByID(id string) (*dtos.CategoryResponse, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	res := &dtos.CategoryResponse{
		ID:  category.ID.String(),
		Name: category.Name,
	}

	return res, nil
}

func (u *categoryUsecase) UpdateCategory(id string, req *dtos.CategoryRequest) (*dtos.CategoryResponse, error) {
	category, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}

	category.Name = req.Name

	if err := u.categoryRepo.Update(category); err != nil {
		return nil, errors.New("failed to update category")
	}

	res := &dtos.CategoryResponse{
		ID: category.ID.String(),
		Name: category.Name,
	}

	return res, nil
}

func (u *categoryUsecase) DeleteCategory(id string) error {
	_, err := u.categoryRepo.FindByID(id)
	if err != nil {
		return errors.New("category not found")
	}

	if err := u.categoryRepo.Delete(id); err != nil {
		return errors.New("failed to delete category")
	}

	return nil
}
