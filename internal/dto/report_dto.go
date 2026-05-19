package dtos

type DailyReportResponse struct {
	Date string `json:"date"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalRevenue float64 `json:"total_revenue"`
}

type MonthlyReportResponse struct {
	Month string `json:"month"`
	TotalTransactions int64 `json:"total_transactions"`
	TotalRevenue float64 `json:"total_revenue"`
}

type RevenueTrendItem struct {
	Date         string  `json:"date"`
	TotalRevenue float64 `json:"total_revenue"`
}
