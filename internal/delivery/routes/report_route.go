package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/database"
	"github.com/savanyv/zenith-pay/internal/delivery/handlers"
	"github.com/savanyv/zenith-pay/internal/middlewares"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func reportRegisterRoutes(app fiber.Router, jwtService helpers.JWTService) {
	repo := repository.NewReportRepository(database.DB)
	uc := usecase.NewReportUsecase(repo)
	handler := handlers.NewReportHandler(uc)

	reportRoutes := app.Group("/admin/reports",
		middlewares.JWTMiddleware(jwtService),
		middlewares.RoleMiddleware(model.AdminRole),
		middlewares.RateLimiter(60, 1*time.Minute),
	)

	reportRoutes.Get("/daily", handler.GetDailyReport)
	reportRoutes.Get("/monthly", handler.GetMonthlyReport)
	reportRoutes.Get("/revenue", handler.GetRevenueTrend)
}
