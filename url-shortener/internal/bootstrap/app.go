package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"shortener/internal/config"
	"shortener/internal/infra/database"
	"shortener/internal/repo"
	urlredis "shortener/internal/repo/url/redis"
	"shortener/internal/service"
	"shortener/internal/transport/http"
)

type Bootstrap struct {
	log *slog.Logger
	cfg *config.Config
}

func (b *Bootstrap) Run() {
	repos, closeDB := b.initRepo()
	defer closeDB()

	svc := service.NewServiceURL(b.cfg.URLShortener, repos, b.log)

	srv := http.NewRestAPI(&b.cfg.HTTP, svc, b.log)

	closeSrv := srv.CreateServer()
	defer closeSrv()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	b.log.Info("shutting down...")
}

func (b *Bootstrap) initRepo() (*repo.Repo, func()) {
	switch b.cfg.Storage.Type {
	case "redis":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		rdb, err := database.GetRedisConnectionPool(ctx, b.cfg.Storage.URL, b.log)

		cancel()

		if err != nil {
			b.log.Error("failed to connect to redis", "error", err)
			return nil, func() {}
		}

		b.log.Info("successful connect to redis")

		urlDB := urlredis.NewURLRedis(rdb, b.log)

		repos := repo.NewRepo(urlDB)

		return repos, func() { _ = rdb.Close() }
	default:
		b.log.Error("unsupported storage type", "type", b.cfg.Storage.Type)
		os.Exit(1)

		return nil, func() {}
	}
}

func NewBootstrap(cfg *config.Config, log *slog.Logger) *Bootstrap {
	return &Bootstrap{
		log: log,
		cfg: cfg,
	}
}
