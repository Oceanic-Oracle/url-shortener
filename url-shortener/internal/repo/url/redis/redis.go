package urlredis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"shortener/internal/repo/url"
	"shortener/internal/service"
)

const urlPrefix = "url:"

type URLRedis struct {
	conn *redis.Client
	log  *slog.Logger
}

func (ur *URLRedis) CreateURL(ctx context.Context, longURL url.LongURL,
	lenShortURL int, ttl time.Duration) (url.ShortURL, error) {
	for i := 0; i < 5; i++ {
		shortURL := service.GenerateShortCode(lenShortURL)

		_, err := ur.conn.Get(ctx, urlPrefix+shortURL).Result()
		if err == nil {
			continue
		}

		if !errors.Is(err, redis.Nil) {
			return "", service.WrapErrInternalServer(fmt.Errorf("redis: %w", err))
		}

		if err = ur.conn.Set(ctx, urlPrefix+shortURL, longURL, ttl*time.Minute).Err(); err != nil {
			return "", service.WrapErrInternalServer(fmt.Errorf("redis: %w", err))
		}
	}

	return "", service.WrapErrInternalServer(errors.New("failed to generate unique short code after 5 attempts"))
}

func (ur *URLRedis) GetURL(ctx context.Context, shortURL url.ShortURL) (url.LongURL, error) {
	return "", nil
}

func NewURLRedis(conn *redis.Client, log *slog.Logger) url.URLInterface {
	return &URLRedis{
		conn: conn,
		log:  log,
	}
}
