package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shortener/internal/config"
	"shortener/internal/metrics"
	"shortener/internal/repo"
	"shortener/internal/service"
	"shortener/internal/transport/http"
)

type Bootstrap struct {
	log *slog.Logger
	cfg *config.Config
}

func (b *Bootstrap) Run() error {
	factory := repo.NewFactory(&b.cfg.Storage, b.log)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repos, closeDB, err := factory.Create(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize repositories: %w", err)
	}
	defer closeDB()

	closeMt := metrics.CreateServer(b.cfg.HTTP, b.log)
	defer closeMt()

	svc := service.NewServiceURL(b.cfg.URLShortener, repos, b.log)

	srv := http.NewRestAPI(&b.cfg.HTTP, svc, b.log)

	closeSrv := srv.CreateServer()
	defer closeSrv()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	b.log.Info("shutting down...")

	return nil
}

func NewBootstrap(cfg *config.Config, log *slog.Logger) *Bootstrap {
	return &Bootstrap{
		log: log,
		cfg: cfg,
	}
}
