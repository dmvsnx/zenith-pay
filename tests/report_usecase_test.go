package tests

import (
	"errors"
	"testing"
	"time"

	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/tests/mocks"
)

func TestReportUsecase_GetDailyReport_Success(t *testing.T) {
	repo := &mocks.ReportRepo{
		GetSalesSummaryFn: func(start, end time.Time) (*repository.SalesSumary, error) {
			return &repository.SalesSumary{
				TotalTransactions: 10,
				TotalRevenue:      500000,
			}, nil
		},
	}
	uc := usecase.NewReportUsecase(repo)

	res, err := uc.GetDailyReport("2025-01-15")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Date != "2025-01-15" {
		t.Fatalf("expected 2025-01-15, got %s", res.Date)
	}
	if res.TotalTransactions != 10 {
		t.Fatalf("expected 10, got %d", res.TotalTransactions)
	}
	if res.TotalRevenue != 500000 {
		t.Fatalf("expected 500000, got %f", res.TotalRevenue)
	}
}

func TestReportUsecase_GetDailyReport_InvalidDate(t *testing.T) {
	uc := usecase.NewReportUsecase(&mocks.ReportRepo{})

	_, err := uc.GetDailyReport("15-01-2025")

	if err == nil || err.Error() != "invalid date format(YYYY-MM-DD)" {
		t.Fatalf("expected invalid date format, got %v", err)
	}
}

func TestReportUsecase_GetDailyReport_RepoError(t *testing.T) {
	repo := &mocks.ReportRepo{
		GetSalesSummaryFn: func(start, end time.Time) (*repository.SalesSumary, error) {
			return nil, errors.New("db error")
		},
	}
	uc := usecase.NewReportUsecase(repo)

	_, err := uc.GetDailyReport("2025-01-15")

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}

func TestReportUsecase_GetMonthlyReport_Success(t *testing.T) {
	repo := &mocks.ReportRepo{
		GetSalesSummaryFn: func(start, end time.Time) (*repository.SalesSumary, error) {
			return &repository.SalesSumary{
				TotalTransactions: 50,
				TotalRevenue:      2500000,
			}, nil
		},
	}
	uc := usecase.NewReportUsecase(repo)

	res, err := uc.GetMonthlyReport("2025-01")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Month != "2025-01" {
		t.Fatalf("expected 2025-01, got %s", res.Month)
	}
	if res.TotalTransactions != 50 {
		t.Fatalf("expected 50, got %d", res.TotalTransactions)
	}
	if res.TotalRevenue != 2500000 {
		t.Fatalf("expected 2500000, got %f", res.TotalRevenue)
	}
}

func TestReportUsecase_GetMonthlyReport_InvalidMonth(t *testing.T) {
	uc := usecase.NewReportUsecase(&mocks.ReportRepo{})

	_, err := uc.GetMonthlyReport("01-2025")

	if err == nil || err.Error() != "invalid month format(YYYY-MM)" {
		t.Fatalf("expected invalid month format, got %v", err)
	}
}

func TestReportUsecase_GetMonthlyReport_RepoError(t *testing.T) {
	repo := &mocks.ReportRepo{
		GetSalesSummaryFn: func(start, end time.Time) (*repository.SalesSumary, error) {
			return nil, errors.New("db error")
		},
	}
	uc := usecase.NewReportUsecase(repo)

	_, err := uc.GetMonthlyReport("2025-01")

	if err == nil || err.Error() != "db error" {
		t.Fatalf("expected db error, got %v", err)
	}
}
