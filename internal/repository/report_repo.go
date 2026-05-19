package repository

import (
	"time"

	"github.com/savanyv/zenith-pay/internal/model"
	"gorm.io/gorm"
)

type SalesSumary struct {
	TotalTransactions int64
	TotalRevenue float64
}

type DailySalesRow struct {
	Date         string
	TotalRevenue float64
}

type ReportRepository interface {
	GetSalesSummary(start, end time.Time) (*SalesSumary, error)
	GetDailySalesInRange(start, end time.Time) ([]*DailySalesRow, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{
		db: db,
	}
}

func (r *reportRepository) GetSalesSummary(start, end time.Time) (*SalesSumary, error) {
	var res SalesSumary

	err := r.db.Model(&model.Transaction{}).
		Select("COUNT(id) as total_transactions, COALESCE(SUM(total_amount), 0) as total_revenue").
		Where("created_at BETWEEN ? AND ?", start, end).
		Scan(&res).Error

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (r *reportRepository) GetDailySalesInRange(start, end time.Time) ([]*DailySalesRow, error) {
	var rows []*DailySalesRow

	err := r.db.Model(&model.Transaction{}).
		Select("DATE(created_at) as date, COALESCE(SUM(total_amount), 0) as total_revenue").
		Where("created_at BETWEEN ? AND ?", start, end).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}
