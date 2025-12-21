package metrics

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"shortener/internal/config"
)

func CreateServer(cfg config.HTTP, log *slog.Logger) func() {
	router := chi.NewRouter()

	router.Handle("/metrics", promhttp.Handler())

	srv := &http.Server{
		Addr:         cfg.AddrMetrics,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		log.Info("HTTP metrics server starting", slog.String("addr", cfg.AddrMetrics))

		if err := srv.ListenAndServe(); err != nil {
			log.Error("HTTP metrics server failed", slog.Any("error", err))
			return
		}
	}()

	return func() {
		if err := srv.Close(); err != nil {
			log.Error("failed to close metrics server", slog.Any("err", err))
		}
	}
}
