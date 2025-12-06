package logctx

import (
	"context"

	"github.com/google/uuid"
)

type Ctx struct {
	ReqID string
}

type keyType int

const key keyType = 0

func WithReqID(ctx context.Context) context.Context {
	reqID := uuid.New().String()

	if c, ok := ctx.Value(key).(Ctx); ok {
		c.ReqID = reqID
		return context.WithValue(ctx, key, c)
	}

	return context.WithValue(ctx, key, Ctx{ReqID: reqID})
}
