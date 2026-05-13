package mocks

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type ProductRepo struct {
	CreateFn             func(product *model.Product) error
	FindByIDFn           func(id string) (*model.Product, error)
	FindByNameFn         func(name string) (*model.Product, error)
	FindBySKUFn          func(sku string) (*model.Product, error)
	FindAllFn            func() ([]*model.Product, error)
	FindAllPaginatedFn   func(offset, limit int) ([]*model.Product, int64, error)
	UpdateFn             func(product *model.Product) error
	DeleteFn             func(id string) error
	FindByIDForUpdateFn  func(tx *gorm.DB, id string) (*model.Product, error)
	UpdateTxFn           func(tx *gorm.DB, product *model.Product) error
}

func (m *ProductRepo) Create(product *model.Product) error { return m.CreateFn(product) }
func (m *ProductRepo) FindByID(id string) (*model.Product, error) { return m.FindByIDFn(id) }
func (m *ProductRepo) FindByName(name string) (*model.Product, error) { return m.FindByNameFn(name) }
func (m *ProductRepo) FindBySKU(sku string) (*model.Product, error) { return m.FindBySKUFn(sku) }
func (m *ProductRepo) FindAll() ([]*model.Product, error) { return m.FindAllFn() }
func (m *ProductRepo) FindAllPaginated(offset, limit int) ([]*model.Product, int64, error) {
	return m.FindAllPaginatedFn(offset, limit)
}
func (m *ProductRepo) Update(product *model.Product) error { return m.UpdateFn(product) }
func (m *ProductRepo) Delete(id string) error { return m.DeleteFn(id) }

func (m *ProductRepo) FindByIDForUpdate(tx *gorm.DB, id string) (*model.Product, error) {
	return m.FindByIDForUpdateFn(tx, id)
}

func (m *ProductRepo) UpdateTx(tx *gorm.DB, product *model.Product) error {
	return m.UpdateTxFn(tx, product)
}
