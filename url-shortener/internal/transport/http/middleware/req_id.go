package middleware

import (
	"net/http"

	"shortener/internal/logctx"
)

func MiddlewareReqID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := logctx.WithReqID(r.Context())
		next(w, r.WithContext(ctx))
	}
}
