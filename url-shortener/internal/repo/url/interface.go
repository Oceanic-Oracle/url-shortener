package url

import (
	"context"
	"time"
)

type URLInterface interface {
	SaveURL(ctx context.Context, shortURL ShortURL, longURL LongURL, ttl time.Duration) error
	GetLongURLByShort(ctx context.Context, shortURL ShortURL) (LongURL, error)
	GetLongURLByShortWithTTLUpdate(ctx context.Context, shortURL ShortURL, ttl time.Duration) (LongURL, error)
}
