package repository

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProductRepository interface {
	Create(product *model.Product) error
	FindByID(id string) (*model.Product, error)
	FindBySKU(sku string) (*model.Product, error)
	FindByName(name string) (*model.Product, error)
	FindAll() ([]*model.Product, error)
	FindAllPaginated(offset, limit int) ([]*model.Product, int64, error)
	Update(product *model.Product) error
	Delete(id string) error
	FindByIDForUpdate(tx *gorm.DB, id string) (*model.Product, error)
	UpdateTx(tx *gorm.DB, product *model.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *model.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").Where("sku = ?", sku).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindByName(name string) (*model.Product, error) {
	var product model.Product
	if err := r.db.Preload("Category").Where("name = ?", name).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) FindAll() ([]*model.Product, error) {
	var products []*model.Product
	if err := r.db.Preload("Category").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *productRepository) FindAllPaginated(offset, limit int) ([]*model.Product, int64, error) {
	var products []*model.Product
	var total int64

	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Category").Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) Update(product *model.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) Delete(id string) error {
	if err := r.db.Delete(&model.Product{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}

func (r *productRepository) FindByIDForUpdate(tx *gorm.DB, id string) (*model.Product, error) {
	var product model.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&product).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) UpdateTx(tx *gorm.DB, product *model.Product) error {
	return tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(map[string]interface{}{
		"stock": product.Stock,
	}).Error
}
