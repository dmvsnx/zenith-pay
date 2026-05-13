package mocks

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type TransactionItemRepo struct {
	CreateManyFn          func(tx *gorm.DB, items []model.TransactionItems) error
	FindByTransactionIDFn func(transactionID string) ([]model.TransactionItems, error)
}

func (m *TransactionItemRepo) CreateMany(tx *gorm.DB, items []model.TransactionItems) error {
	return m.CreateManyFn(tx, items)
}

func (m *TransactionItemRepo) FindByTransactionID(transactionID string) ([]model.TransactionItems, error) {
	return m.FindByTransactionIDFn(transactionID)
}
