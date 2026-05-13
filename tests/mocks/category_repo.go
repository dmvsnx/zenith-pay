package mocks

import "github.com/savanyv/zenith-pay/internal/model"

type CategoryRepo struct {
	CreateFn             func(category *model.Category) error
	FindByNameFn         func(name string) (*model.Category, error)
	FindByIDFn           func(id string) (*model.Category, error)
	FindAllFn            func() ([]*model.Category, error)
	FindAllPaginatedFn   func(offset, limit int) ([]*model.Category, int64, error)
	UpdateFn             func(category *model.Category) error
	DeleteFn             func(id string) error
}

func (m *CategoryRepo) Create(category *model.Category) error {
	return m.CreateFn(category)
}

func (m *CategoryRepo) FindByName(name string) (*model.Category, error) {
	return m.FindByNameFn(name)
}

func (m *CategoryRepo) FindByID(id string) (*model.Category, error) {
	return m.FindByIDFn(id)
}

func (m *CategoryRepo) FindAll() ([]*model.Category, error) {
	return m.FindAllFn()
}

func (m *CategoryRepo) FindAllPaginated(offset, limit int) ([]*model.Category, int64, error) {
	return m.FindAllPaginatedFn(offset, limit)
}

func (m *CategoryRepo) Update(category *model.Category) error {
	return m.UpdateFn(category)
}

func (m *CategoryRepo) Delete(id string) error {
	return m.DeleteFn(id)
}
