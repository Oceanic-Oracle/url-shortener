package middleware

import (
	"net/http"

	"shortener/internal/logctx"
)

func MiddlewareReqID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
