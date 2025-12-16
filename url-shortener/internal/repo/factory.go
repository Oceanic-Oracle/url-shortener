package repo

import (
	"context"
	"errors"
	"log/slog"

	"shortener/internal/config"
	"shortener/internal/infra/database"
	urlredis "shortener/internal/repo/url/redis"
)

var ErrUnsupportedStorage = errors.New("unsupported storage type")

type RepoFactory struct {
	cfg *config.Storage
	log *slog.Logger
}

func (f *RepoFactory) Create(ctx context.Context) (*Repo, func(), error) {
	switch f.cfg.Type {
	case "redis":
		return f.setupRedis(ctx)
	default:
		f.log.Error("unsupported storage type", "type", f.cfg.Type)

		return nil, func() {}, ErrUnsupportedStorage
	}
}

func (f *RepoFactory) setupRedis(ctx context.Context) (*Repo, func(), error) {
	rdb, err := database.GetRedisConnectionPool(ctx, f.cfg.URL, f.log)
	if err != nil {
		f.log.Error("failed to connect to redis", slog.Any("err", err))
		return nil, func() {}, err
	}

	f.log.Info("successful connect to redis")

	urlDB := urlredis.NewURLRedis(rdb, f.log)

	repos := NewRepo(urlDB)

	return repos, func() {
		if err := rdb.Close(); err != nil {
			f.log.Error("failed to close Redis", slog.Any("err", err))
		}
	}, nil
}

func NewFactory(cfg *config.Storage, log *slog.Logger) *RepoFactory {
	return &RepoFactory{
		cfg: cfg,
		log: log,
	}
}
