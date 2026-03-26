package usecase

import (
	"errors"
	"time"

	dtos "github.com/savanyv/zenith-pay/internal/dto"
	"github.com/savanyv/zenith-pay/internal/repository"
)

type ReportUsecase interface {
	GetDailyReport(date string) (*dtos.DailyReportResponse, error)
	GetMonthlyReport(month string) (*dtos.MonthlyReportResponse, error)
}

type reportUsecase struct {
	reportRepo repository.ReportRepository
}

func NewReportUsecase(rr repository.ReportRepository) ReportUsecase {
	return &reportUsecase{
		reportRepo: rr,
	}
}

func (u *reportUsecase) GetDailyReport(date string) (*dtos.DailyReportResponse, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	parsedDate, err := time.ParseInLocation("2006-01-02", date, loc)
	if err != nil {
		return nil, errors.New("invalid date format(YYYY-MM-DD)")
	}

	start := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, loc)
	end := start.Add(24 * time.Hour - time.Nanosecond)

	summary, err := u.reportRepo.GetSalesSummary(start, end)
	if err != nil {
		return nil, err
	}

	res := &dtos.DailyReportResponse{
		Date: start.Format("2006-01-02"),
		TotalTransactions: summary.TotalTransactions,
		TotalRevenue: summary.TotalRevenue,
	}

	return res, nil
}

func (u *reportUsecase) GetMonthlyReport(month string) (*dtos.MonthlyReportResponse, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	parsedMonth, err := time.ParseInLocation("2006-01", month, loc)
	if err != nil {
		return nil, errors.New("invalid month format(YYYY-MM)")
	}

	start := time.Date(parsedMonth.Year(), parsedMonth.Month(), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0).Add(-time.Nanosecond)

	summary, err := u.reportRepo.GetSalesSummary(start, end)
	if err != nil {
		return nil, err
	}

	res := &dtos.MonthlyReportResponse{
		Month: start.Format("2006-01"),
		TotalTransactions: summary.TotalTransactions,
		TotalRevenue: summary.TotalRevenue,
	}

	return res, nil
}
