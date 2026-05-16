package model

import (
	"time"

	"github.com/google/uuid"
)

type PaymentMethod string

const (
	Cash  PaymentMethod = "cash"
	Debit PaymentMethod = "debit"
	Qris  PaymentMethod = "qris"
)

func (p PaymentMethod) IsValid() bool {
	switch p {
	case Cash, Debit, Qris:
		return true
	default:
		return false
	}
}

type Transaction struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	UserID uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	User   User      `json:"user" gorm:"foreignKey:UserID;references:ID"`

	ShiftID uuid.UUID `json:"shift_id" gorm:"type:uuid;not null"`
	Shift   Shift     `json:"shift" gorm:"foreignKey:ShiftID;references:ID"`

	TransactionDate time.Time     `json:"transaction_date" gorm:"not null"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"type:varchar(20);not null"`

	TotalAmount   int64 `json:"total_amount" gorm:"not null"`
	PaymentAmount int64 `json:"payment_amount" gorm:"not null"`
	ChangeAmount  int64 `json:"change_amount" gorm:"not null"`

	TransactionItems []TransactionItems `json:"transaction_items" gorm:"foreignKey:TransactionID;references:ID"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
