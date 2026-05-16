package dtos

import "time"

// ===== REQUEST =====

type TransactionItemRequest struct {
	ProductID string `json:"product_id" validate:"required,uuid"`
	Quantity  int    `json:"quantity" validate:"required,gt=0"`
}

type TransactionRequest struct {
	PaymentMethod string                   `json:"payment_method" validate:"required,oneof=cash debit qris"`
	PaymentAmount int64                    `json:"payment_amount"`
	Items         []TransactionItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ===== RESPONSE =====

type TransactionItemResponse struct {
	ProductID    string `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice int64  `json:"product_price"`
	Quantity     int    `json:"quantity"`
	SubTotal     int64  `json:"sub_total"`
}

type TransactionResponse struct {
	ID              string                    `json:"id"`
	UserID          string                    `json:"user_id"`
	TransactionDate time.Time                 `json:"transaction_date"`
	PaymentMethod   string                    `json:"payment_method"`
	TotalAmount     int64                     `json:"total_amount"`
	PaymentAmount   int64                     `json:"payment_amount"`
	ChangeAmount    int64                     `json:"change_amount"`
	Items           []TransactionItemResponse `json:"items,omitempty"`
}
