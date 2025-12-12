package main

import (
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

	app.NewBootstrap(cfg, log).Run()
}
