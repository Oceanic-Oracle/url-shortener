package middleware

import (
	"net/http"
	"time"
	"unicode"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path"},
	)

	requestDuration = promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_seconds",
			Help:       "HTTP request duration in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"method", "path"},
	)
)

func ObserveMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		httpRequestsTotal.WithLabelValues(r.Method, normalizePath(r.URL.Path)).Inc()
		requestDuration.WithLabelValues(r.Method, normalizePath(r.URL.Path)).Observe(time.Since(start).Seconds())
	})
}

func normalizePath(path string) string {
	for _, r := range path[1:] {
		if unicode.IsUpper(r) || unicode.IsDigit(r) {
			return "/{code}"
		} else if r == '/' || r == '-' {
			return path
		}
	}

	return path
}
