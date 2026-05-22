package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/config"
	cld "github.com/savanyv/zenith-pay/internal/utils/cloudinary"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RegisterRoutes(app fiber.Router) {
	cfg := config.LoadConfig()

	jwtService := helpers.NewJWTService()
	bcrypt := helpers.NewBcryptHelper()

	cloudinaryService, err := cld.NewCloudinaryService(cfg.CloudinaryURL)
	if err != nil {
		cloudinaryService = nil
	}

	healthRegisterRoutes(app)

	api := app.Group("/zenith-pay")

	userRegisterRoutes(api, jwtService, bcrypt)
	categoryRegisterRoutes(api, jwtService)
	productRegisterRoutes(api, jwtService, cloudinaryService)
	transactionRegisterRoutes(api, jwtService)
	shiftRegisterRoute(api, jwtService)
	reportRegisterRoutes(api, jwtService)
}
