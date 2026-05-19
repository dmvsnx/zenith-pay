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

func transactionRegisterRoutes(app fiber.Router, jwtService helpers.JWTService) {
	repo := repository.NewTransactionRepository(database.DB)
	itemRepo := repository.NewTransactionItemRepository(database.DB)
	productRepo := repository.NewProductRepository(database.DB)
	shiftRepo := repository.NewShiftRepository(database.DB)
	uc := usecase.NewTransactionUsecase(
		database.DB,
		repo,
		itemRepo,
		productRepo,
	)
	handler := handlers.NewTransactionHandler(uc)

	transactionRoutes := app.Group(
		"/transactions",
		middlewares.JWTMiddleware(jwtService),
		middlewares.RoleMiddleware(model.CashierRole),
		middlewares.RateLimiter(30, 1*time.Minute),
	)
	transactionRoutes.Post("/", middlewares.RequireActiveShift(shiftRepo), handler.CreateTransaction)
	transactionRoutes.Get("/", handler.ListMyTransactions)
	transactionRoutes.Get("/:id", handler.GetMyTransactionByID)

	adminRoutes := app.Group(
		"/admin/transactions",
		middlewares.JWTMiddleware(jwtService),
		middlewares.RoleMiddleware(model.AdminRole),
		middlewares.RateLimiter(60, 1*time.Minute),
	)
	adminRoutes.Get("/", handler.ListTransactions)
	adminRoutes.Get("/:id", handler.GetTransactionByID)
}
