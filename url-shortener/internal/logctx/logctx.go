package logctx

import (
	"context"
	"log/slog"
)

type HandlerMiddleware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddleware {
	return &HandlerMiddleware{next: next}
}

func (h *HandlerMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *HandlerMiddleware) Handle(ctx context.Context, record slog.Record) error {
	if c, ok := ctx.Value(key).(Ctx); ok {
		if c.ReqID != "" {
			record.Add("req_id", c.ReqID)
		}

		if c.Code != "" {
			record.Add("code", c.Code)
		}

		if c.URL != "" {
			record.Add("url", c.URL)
		}
	}

	return h.next.Handle(ctx, record)
}

func (h *HandlerMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewHandlerMiddleware(h.next.WithAttrs(attrs))
}

func (h *HandlerMiddleware) WithGroup(name string) slog.Handler {
	return NewHandlerMiddleware(h.next.WithGroup(name))
}
