package repository

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *model.Category) error
	FindByName(name string) (*model.Category, error)
	FindByID(id string) (*model.Category, error)
	FindAll() ([]*model.Category, error)
	FindAllPaginated(offset, limit int) ([]*model.Category, int64, error)
	Update(category *model.Category) error
	Delete(id string) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(category *model.Category) error {
	if err := r.db.Create(category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) FindByName(name string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("name = ?", name).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindByID(id string) (*model.Category, error) {
	var category model.Category
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) FindAll() ([]*model.Category, error) {
	var categories []*model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindAllPaginated(offset, limit int) ([]*model.Category, int64, error) {
	var categories []*model.Category
	var total int64

	if err := r.db.Model(&model.Category{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) Update(category *model.Category) error {
	if err := r.db.Save(category).Error; err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) Delete(id string) error {
	if err := r.db.Delete(&model.Category{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
