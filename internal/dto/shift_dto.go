package dtos

import "time"

type OpenShiftRequest struct {
	OpeningBalance int64 `json:"opening_balance" validate:"required,gt=0"`
}

type CloseShiftRequest struct {
	ShiftID string `json:"shift_id" validate:"required,uuid"`
	ClosingBalance int64 `json:"closing_balance" validate:"required,gt=0"`
}

type ShiftResponse struct {
	ID              string     `json:"id"`
	CashierID       string     `json:"cashier_id"`
	Status          string     `json:"status"`
	OpeningBalance  int64      `json:"opening_balance"`
	ClosingBalance  *int64     `json:"closing_balance,omitempty"`
	ExpectedClosingBalance *int64 `json:"expected_closing_balance,omitempty"`
	Variance        *int64     `json:"variance,omitempty"`
	OpenedAt        time.Time  `json:"opened_at"`
	ClosedAt        *time.Time `json:"closed_at,omitempty"`
}

