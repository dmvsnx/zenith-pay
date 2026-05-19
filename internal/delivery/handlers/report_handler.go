package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

type ReportHandler struct {
	reportUsecase usecase.ReportUsecase
}

func NewReportHandler(ru usecase.ReportUsecase) *ReportHandler {
	return &ReportHandler{
		reportUsecase: ru,
	}
}

func (h *ReportHandler) GetDailyReport(c *fiber.Ctx) error {
	date := c.Query("date")
	if date == "" {
		return helpers.BadRequest(c, "Date parameter is required(YYYY-MM-DD)")
	}

	if _, err := time.Parse("2006-01-02", date); err != nil {
		return helpers.BadRequest(c, "Invalid date format(YYYY-MM-DD)")
	}

	res, err := h.reportUsecase.GetDailyReport(date)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Daily report retrieved successfully", res)
}

func (h *ReportHandler) GetRevenueTrend(c *fiber.Ctx) error {
	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {
		return helpers.BadRequest(c, "from and to parameters are required(YYYY-MM-DD)")
	}

	if _, err := time.Parse("2006-01-02", from); err != nil {
		return helpers.BadRequest(c, "Invalid from date format(YYYY-MM-DD)")
	}

	if _, err := time.Parse("2006-01-02", to); err != nil {
		return helpers.BadRequest(c, "Invalid to date format(YYYY-MM-DD)")
	}

	res, err := h.reportUsecase.GetRevenueTrend(from, to)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Revenue trend retrieved successfully", res)
}

func (h *ReportHandler) GetMonthlyReport(c *fiber.Ctx) error {
	period := c.Query("period")
	if period == "" {
		return helpers.BadRequest(c, "Period parameter is required(YYYY-MM)")
	}

	if _, err := time.Parse("2006-01", period); err != nil {
		return helpers.BadRequest(c, "Invalid period format(YYYY-MM)")
	}

	res, err := h.reportUsecase.GetMonthlyReport(period)
	if err != nil {
		return helpers.InternalServerError(c, err.Error())
	}

	return helpers.Success(c, "Monthly report retrieved successfully", res)
}
