package mocks

import (
	"time"

	"github.com/savanyv/zenith-pay/internal/repository"
)

type ReportRepo struct {
	GetSalesSummaryFn    func(start, end time.Time) (*repository.SalesSumary, error)
	GetDailySalesInRangeFn func(start, end time.Time) ([]*repository.DailySalesRow, error)
}

func (m *ReportRepo) GetSalesSummary(start, end time.Time) (*repository.SalesSumary, error) {
	return m.GetSalesSummaryFn(start, end)
}

func (m *ReportRepo) GetDailySalesInRange(start, end time.Time) ([]*repository.DailySalesRow, error) {
	return m.GetDailySalesInRangeFn(start, end)
}
