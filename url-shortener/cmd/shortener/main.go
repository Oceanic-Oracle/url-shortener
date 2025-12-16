package main

import (
	"os"

	"shortener/internal/bootstrap"
	"shortener/internal/config"
	"shortener/internal/infra/logger"
)

const (
	PATH = "../.env"
)

func main() {
	cfg := config.MustLoad(PATH)

	log := logger.SetupLogger(cfg.Env)

	if err := app.NewBootstrap(cfg, log).Run(); err != nil {
		log.Error("application failed", "error", err)
		os.Exit(1)
	}

	log.Info("application exited successfully")
}
