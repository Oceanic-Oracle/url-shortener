package url

import (
	"context"
	"time"
)

type URLInterface interface {
	SaveURL(ctx context.Context, shortURL ShortURL, longURL LongURL, ttl time.Duration) error
	GetShortURLByLong(ctx context.Context, longURL LongURL) (ShortURL, error)
	GetLongURLByShort(ctx context.Context, shortURL ShortURL) (LongURL, error)
	GetShortURLByLongWithTTLUpdate(ctx context.Context, longURL LongURL, ttl time.Duration) (ShortURL, error)
	GetLongURLByShortWithTTLUpdate(ctx context.Context, shortURL ShortURL, ttl time.Duration) (LongURL, error)
}
