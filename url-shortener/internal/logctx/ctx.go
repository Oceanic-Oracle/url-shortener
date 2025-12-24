package logctx

import (
	"context"

	"github.com/google/uuid"
)

type Ctx struct {
	ReqID  string
	URL    string
	Code   string
	Status int
}

type keyType int

const key keyType = 0

func WithReqID(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, Ctx{ReqID: uuid.New().String()})
}

func WithURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, key, Ctx{URL: url})
}

func WithCode(ctx context.Context, code string) context.Context {
	return context.WithValue(ctx, key, Ctx{Code: code})
}

func GetReqId(ctx context.Context) string {
	if c, ok := ctx.Value(key).(Ctx); ok && c.ReqID != "" {
		return c.ReqID
	}

	return "unknown"
}
