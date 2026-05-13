package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/delivery/handlers"
)

func healthRegisterRoutes(app fiber.Router) {
	app.Get("/health", handlers.LivenessCheck)
	app.Get("/health/ready", handlers.ReadinessCheck)
}
