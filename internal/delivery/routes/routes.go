package routes

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/config"
	"github.com/savanyv/zenith-pay/internal/storage/minio"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

func RegisterRoutes(app fiber.Router) {
	cfg := config.LoadConfig()

	jwtService := helpers.NewJWTService()
	bcrypt := helpers.NewBcryptHelper()

	minioClient, err := minio.New(cfg.MinioEndpoint, cfg.MinioAccessKey, cfg.MinioSecretKey, cfg.MinioBucket, false)
	if err != nil {
		log.Printf("Warning: failed to connect to MinIO: %v", err)
	}
	var minioService minio.Service = minioClient

	healthRegisterRoutes(app)

	api := app.Group("/zenith-pay")

	userRegisterRoutes(api, jwtService, bcrypt)
	categoryRegisterRoutes(api, jwtService)
	productRegisterRoutes(api, jwtService, minioService)
	transactionRegisterRoutes(api, jwtService)
	shiftRegisterRoute(api, jwtService)
	reportRegisterRoutes(api, jwtService)
}
