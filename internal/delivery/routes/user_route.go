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

func userRegisterRoutes(app fiber.Router, jwtService helpers.JWTService, bcrypt helpers.BcryptHelper) {
	repo := repository.NewUserRepository(database.DB)
	usecase := usecase.NewUserUsecase(repo, jwtService, bcrypt)
	handler := handlers.NewUserHandler(usecase)

	auth := app.Group("/auth")
	auth.Post("/login", middlewares.RateLimiter(5, 1*time.Minute),handler.Login)

	admin := app.Group("/admin/users", middlewares.JWTMiddleware(jwtService), middlewares.RoleMiddleware(model.AdminRole), middlewares.RateLimiter(20, 1*time.Minute))
	admin.Get("/", handler.ListUsers)
	admin.Post("/", handler.Register)
}
