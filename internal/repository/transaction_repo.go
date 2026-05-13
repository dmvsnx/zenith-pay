package repository

import (
	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *model.Transaction) error
	FindByID(id string) (*model.Transaction, error)
	FindAll() ([]*model.Transaction, error)
	FindAllPaginated(offset, limit int) ([]*model.Transaction, int64, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(tx *gorm.DB, transaction *model.Transaction) error {
	db := r.db
	if tx != nil {
		db = tx
	}
	return db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.db.Preload("User").Preload("TransactionItems").Where("id = ?", id).First(&transaction).Error; err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *transactionRepository) FindAll() ([]*model.Transaction, error) {
      var transactions []*model.Transaction
      if err := r.db.Preload("TransactionItems").Order("created_at desc").Find(&transactions).Error; err != nil {
            return nil, err
      }
      return transactions, nil
}

func (r *transactionRepository) FindAllPaginated(offset, limit int) ([]*model.Transaction, int64, error) {
      var transactions []*model.Transaction
      var total int64

      if err := r.db.Model(&model.Transaction{}).Count(&total).Error; err != nil {
            return nil, 0, err
      }

      if err := r.db.Preload("TransactionItems").Order("created_at desc").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
            return nil, 0, err
      }

      return transactions, total, nil
}
