package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/savanyv/zenith-pay/config"
	"github.com/savanyv/zenith-pay/internal/database"
	"github.com/savanyv/zenith-pay/internal/database/seed"
	"github.com/savanyv/zenith-pay/internal/delivery/routes"
	"github.com/savanyv/zenith-pay/internal/middlewares"
	"github.com/savanyv/zenith-pay/internal/utils/helpers"
	"github.com/savanyv/zenith-pay/internal/utils/logger"
)

type Server struct {
	app    *fiber.App
	config *config.Config
}

func NewServer(config *config.Config) *Server {
	app := fiber.New(fiber.Config{
		AppName:      config.AppName,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	})

	return &Server{
		app:    app,
		config: config,
	}
}

func (s *Server) Start() error {
	if _, err := database.InitDatabase(s.config); err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	if s.config.AppEnv == "development" {
		bcHelper := helpers.NewBcryptHelper()
		seed.SeedAdmin(database.DB, bcHelper)
	}

	s.app.Use(recover.New())
	s.app.Use(requestid.New())
	s.app.Use(middlewares.CORSMiddleware())
	s.app.Use(middlewares.MethodValidationMiddleware())
	s.app.Use(fiberlogger.New(fiberlogger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	routes.RegisterRoutes(s.app)

	addr := fmt.Sprintf(":%s", s.config.AppPort)
	go func() {
		logger.Log.Info().Str("addr", addr).Msg("Server running")
		if err := s.app.Listen(addr); err != nil {
			logger.Log.Error().Err(err).Msg("Server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.app.ShutdownWithContext(ctx); err != nil {
		return fmt.Errorf("shutdown server: %w", err)
	}

	logger.Log.Info().Msg("Server shut down successfully")
	return nil
}
