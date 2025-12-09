package database

import (
	"context"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

func GetRedisConnectionPool(ctx context.Context, url string, log *slog.Logger) (*redis.Client, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Error("parse redis url", "err", err)
		return nil, err
	}

	rdb := redis.NewClient(opt)

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Error("failed to connect to redis", "err", err)
		return nil, err
	}

	return rdb, nil
}
