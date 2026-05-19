package mocks

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type TransactionRepo struct {
	CreateFn             func(tx *gorm.DB, transaction *model.Transaction) error
	FindByIDFn           func(id string) (*model.Transaction, error)
	FindAllFn            func() ([]*model.Transaction, error)
	FindAllPaginatedFn   func(offset, limit int, from, to string, userID *string) ([]*model.Transaction, int64, error)
	SumByShiftIDGroupedFn func(shiftID string) (cash, debit, qris int64, err error)
}

func (m *TransactionRepo) Create(tx *gorm.DB, transaction *model.Transaction) error {
	return m.CreateFn(tx, transaction)
}
func (m *TransactionRepo) FindByID(id string) (*model.Transaction, error) { return m.FindByIDFn(id) }
func (m *TransactionRepo) FindAll() ([]*model.Transaction, error)         { return m.FindAllFn() }
func (m *TransactionRepo) FindAllPaginated(offset, limit int, from, to string, userID *string) ([]*model.Transaction, int64, error) {
	return m.FindAllPaginatedFn(offset, limit, from, to, userID)
}
func (m *TransactionRepo) SumByShiftIDGrouped(shiftID string) (cash, debit, qris int64, err error) {
	return m.SumByShiftIDGroupedFn(shiftID)
}
