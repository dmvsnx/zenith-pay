package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/savanyv/zenith-pay/internal/database"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
)

var (
	startTime time.Time
	once      sync.Once
)

func init() {
	once.Do(func() {
		startTime = time.Now()
	})
}

func LivenessCheck(c *fiber.Ctx) error {
	return helpers.Success(c, "Service is alive", fiber.Map{
		"service": "zenith-pay",
		"uptime":  time.Since(startTime).String(),
	})
}

func ReadinessCheck(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := database.DB.DB()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(helpers.APIResponse{
			Code:    fiber.StatusServiceUnavailable,
			Status:  "error",
			Message: "Service is not ready",
			Data: fiber.Map{
				"database": "disconnected",
			},
		})
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(helpers.APIResponse{
			Code:    fiber.StatusServiceUnavailable,
			Status:  "error",
			Message: "Service is not ready",
			Data: fiber.Map{
				"database": "unreachable",
			},
		})
	}

	return helpers.Success(c, "Service is ready", fiber.Map{
		"database": "connected",
	})
}
