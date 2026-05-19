package repository

import (
	"time"

	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(tx *gorm.DB, transaction *model.Transaction) error
	FindByID(id string) (*model.Transaction, error)
	FindAll() ([]*model.Transaction, error)
	FindAllPaginated(offset, limit int, from, to string, userID *string) ([]*model.Transaction, int64, error)
	SumByShiftIDGrouped(shiftID string) (cash, debit, qris int64, err error)
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

func (r *transactionRepository) FindAllPaginated(offset, limit int, from, to string, userID *string) ([]*model.Transaction, int64, error) {
	var transactions []*model.Transaction
	var total int64

	query := r.db.Model(&model.Transaction{})
	if from != "" {
		if fromTime, err := time.Parse("2006-01-02", from); err == nil {
			query = query.Where("transaction_date >= ?", fromTime)
		}
	}
	if to != "" {
		if toTime, err := time.Parse("2006-01-02", to); err == nil {
			toTime = toTime.Add(24*time.Hour - time.Second)
			query = query.Where("transaction_date <= ?", toTime)
		}
	}
	if userID != nil && *userID != "" {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("TransactionItems").Order("created_at desc").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *transactionRepository) SumByShiftIDGrouped(shiftID string) (cash, debit, qris int64, err error) {
	type PaymentTotal struct {
		PaymentMethod string
		Total         int64
	}
	var results []PaymentTotal
	err = r.db.Model(&model.Transaction{}).
		Select("payment_method, COALESCE(SUM(total_amount)::bigint, 0) as total").
		Where("shift_id = ?", shiftID).
		Group("payment_method").
		Scan(&results).Error
	if err != nil {
		return 0, 0, 0, err
	}
	for _, r := range results {
		switch r.PaymentMethod {
		case string(model.Cash):
			cash = r.Total
		case string(model.Debit):
			debit = r.Total
		case string(model.Qris):
			qris = r.Total
		}
	}
	return cash, debit, qris, nil
}
