package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RegisterRoutes(app fiber.Router) {
	jwtService := helpers.NewJWTService()
	bcrypt := helpers.NewBcryptHelper()

	api := app.Group("/zenith-pay")

	userRegisterRoutes(api, jwtService, bcrypt)
	categoryRegisterRoutes(api, jwtService)
	productRegisterRoutes(api, jwtService)
	transactionRegisterRoutes(api, jwtService)
	shiftRegisterRoute(api, jwtService)
	reportRegisterRoutes(api, jwtService)
}
