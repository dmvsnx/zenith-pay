package main

import (
	"os"

	"github.com/savanyv/zenith-pay/config"
	"github.com/savanyv/zenith-pay/internal/app"
	"github.com/savanyv/zenith-pay/internal/utils/logger"
)

func main() {
	cfg := config.LoadConfig()
	logger.Init(cfg.AppEnv)

	logger.Log.Info().
		Str("app", cfg.AppName).
		Str("env", cfg.AppEnv).
		Msg("Starting server")

	server := app.NewServer(cfg)
	if err := server.Start(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server stopped with error")
		os.Exit(1)
	}
}
