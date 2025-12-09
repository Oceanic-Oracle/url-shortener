package urlredis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"shortener/internal/repo/url"
)

const urlPrefix = "url:"

type URLRedis struct {
	conn *redis.Client
	log  *slog.Logger
}

func (ur *URLRedis) SaveURL(ctx context.Context, shortURL url.ShortURL, longURL url.LongURL, ttl time.Duration) error {
	ok, err := ur.conn.SetNX(ctx, urlPrefix+string(shortURL), string(longURL), ttl).Result()
	if err != nil {
		return err
	}

	if !ok {
		return url.ErrURLExists
	}

	return nil
}

func (ur *URLRedis) GetLongURLByShort(ctx context.Context, shortURL url.ShortURL) (url.LongURL, error) {
	answ, err := ur.conn.Get(ctx, string(urlPrefix+shortURL)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url.ErrURLNotFound
		}

		return "", err
	}

	return url.LongURL(answ), nil
}

func (ur *URLRedis) GetLongURLByShortWithTTLUpdate(ctx context.Context, shortURL url.ShortURL,
	ttl time.Duration) (url.LongURL, error) {
	answ, err := ur.conn.GetEx(ctx, string(urlPrefix+shortURL), ttl).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url.ErrURLNotFound
		}

		return "", err
	}

	return url.LongURL(answ), nil
}

func NewURLRedis(conn *redis.Client, log *slog.Logger) url.URLInterface {
	return &URLRedis{
		conn: conn,
		log:  log,
	}
}
