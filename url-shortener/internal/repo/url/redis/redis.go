package urlredis

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"

	"shortener/internal/repo/url"
)

const urlPrefixShort = "short:"
const urlPrefixLong = "long:"

type URLRedis struct {
	conn *redis.Client
	log  *slog.Logger
}

func (ur *URLRedis) SaveURL(ctx context.Context, shortURL url.ShortURL, longURL url.LongURL, ttl time.Duration) error {
	ok, err := ur.conn.SetNX(ctx, string(urlPrefixShort+shortURL), longURL, ttl).Result()
	if err != nil {
		return err
	}

	if !ok {
		return url.ErrShortURLExists
	}

	ok, err = ur.conn.SetNX(ctx, urlPrefixLong+string(longURL), string(shortURL), ttl).Result()
	if err != nil {
		ur.conn.Del(ctx, urlPrefixShort+string(shortURL))
		return err
	}

	if !ok {
		ur.conn.Del(ctx, urlPrefixShort+string(shortURL))
		return url.ErrLongURLExists
	}

	return nil
}

func (ur *URLRedis) GetShortURLByLong(ctx context.Context, longURL url.LongURL) (url.ShortURL, error) {
	answ, err := ur.conn.Get(ctx, string(urlPrefixLong+longURL)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url.ErrURLNotFound
		}

		return "", err
	}

	return url.ShortURL(answ), nil
}

func (ur *URLRedis) GetLongURLByShort(ctx context.Context, shortURL url.ShortURL) (url.LongURL, error) {
	answ, err := ur.conn.Get(ctx, string(urlPrefixShort+shortURL)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url.ErrURLNotFound
		}

		return "", err
	}

	return url.LongURL(answ), nil
}

func (ur *URLRedis) GetShortURLByLongWithTTLUpdate(ctx context.Context, longURL url.LongURL,
	ttl time.Duration) (url.ShortURL, error) {
	answ, err := ur.conn.GetEx(ctx, string(urlPrefixLong+longURL), ttl).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", url.ErrURLNotFound
		}

		return "", err
	}

	return url.ShortURL(answ), nil
}

func (ur *URLRedis) GetLongURLByShortWithTTLUpdate(ctx context.Context, shortURL url.ShortURL,
	ttl time.Duration) (url.LongURL, error) {
	answ, err := ur.conn.GetEx(ctx, string(urlPrefixShort+shortURL), ttl).Result()
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
