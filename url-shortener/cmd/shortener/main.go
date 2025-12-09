package main

import (
	"shortener/internal/bootstrap"
	"shortener/internal/config"
	"shortener/internal/infra/logger"
)

func main() {
	cfg := config.MustLoad()

	log := logger.SetupLogger(cfg.Env)

	app.NewBootstrap(cfg, log).Run()
}
