package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/database"
	"github.com/savanyv/zenith-pay/internal/delivery/handlers"
	"github.com/savanyv/zenith-pay/internal/middlewares"
	"github.com/savanyv/zenith-pay/internal/model"
	"github.com/savanyv/zenith-pay/internal/repository"
	"github.com/savanyv/zenith-pay/internal/storage/minio"
	"github.com/savanyv/zenith-pay/internal/usecase"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func productRegisterRoutes(app fiber.Router, jwtService helpers.JWTService, minioService minio.Service) {
	repo := repository.NewProductRepository(database.DB)
	categoryRepo := repository.NewCategoryRepository(database.DB)
	uc := usecase.NewProductUsecase(repo, categoryRepo, minioService)
	handler := handlers.NewProductHandler(uc, minioService)

	productRoutes := app.Group("/products", middlewares.JWTMiddleware(jwtService), middlewares.RateLimiter(100, 1*time.Minute))
	productRoutes.Get("/", handler.ListProduct)
	productRoutes.Get("/:id", handler.GetProductByID)

	admin := productRoutes.Group("/admin", middlewares.RoleMiddleware(model.AdminRole), middlewares.RateLimiter(50, 1*time.Minute))
	admin.Post("/", handler.CreateProduct)
	admin.Put("/:id", handler.UpdateProduct)
	admin.Delete("/:id", handler.DeleteProduct)
}
