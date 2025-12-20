package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"

	"shortener/internal/config"
	"shortener/internal/service"
	"shortener/internal/transport/http/middleware"
	"shortener/internal/transport/http/url"
)

type Server struct {
	log *slog.Logger
	cfg *config.HTTP
	svc *service.ServiceURL
}

func (s *Server) CreateServer() func() {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ping"))
	})

	router.Post("/shorten", middleware.MiddlewareReqID(url.CreateURL(s.svc, s.cfg.Host, s.log)))
	router.Get("/{code}", middleware.MiddlewareReqID(url.RedirectURL(s.svc, s.log)))

	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		router.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      corsHandler,
		ReadTimeout:  s.cfg.Timeout,
		WriteTimeout: s.cfg.Timeout,
		IdleTimeout:  s.cfg.IdleTimeout,
	}

	go func() {
		s.log.Info("HTTP server starting", slog.String("addr", s.cfg.Addr))

		if err := srv.ListenAndServe(); err != nil {
			s.log.Error("HTTP server failed", slog.Any("error", err))
			return
		}
	}()

	return func() {
		if err := srv.Close(); err != nil {
			s.log.Error("failed to close server", slog.Any("err", err))
		}
	}
}

func NewRestAPI(cfg *config.HTTP, svc *service.ServiceURL, log *slog.Logger) *Server {
	return &Server{
		cfg: cfg,
		svc: svc,
		log: log,
	}
}
