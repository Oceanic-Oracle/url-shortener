package middleware

import (
	"net/http"
	"shortener/internal/logctx"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var requestMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "shortener",
	Subsystem:  "http",
	Name:       "request_duration_seconds",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"method", "path", "req_id", "status"})

func ObserveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		ww := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(ww, r)

		reqId := logctx.GetReqId(r.Context())

		requestMetrics.WithLabelValues(r.Method, r.URL.Path, reqId,
			strconv.Itoa(ww.statusCode)).Observe(time.Since(start).Seconds(),
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}
