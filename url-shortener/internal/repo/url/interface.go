package url

import (
	"context"
	"time"
)

type URLInterface interface {
	CreateURL(context.Context, LongURL, int, time.Duration) (ShortURL, error)
	GetURL(context.Context, ShortURL) (LongURL, error)
}
