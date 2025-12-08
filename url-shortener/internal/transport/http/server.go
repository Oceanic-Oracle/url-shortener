package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"

	"shortener/internal/config"
	"shortener/internal/service"
	"shortener/internal/transport/http/url"
)

type Server struct {
	log *slog.Logger
	cfg *config.HTTP
	svc *service.ServiceURL
}

func (s *Server) CreateServer() func() {
	router := chi.NewRouter()

	router.Post("/shorten", url.CreateURL(s.svc, s.log))
	router.Get("/{code}", url.RedirectURL(s.svc, s.log))

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

	idleTimeout, _ := strconv.Atoi(s.cfg.IdleTimeout)
	timeout, _ := strconv.Atoi(s.cfg.Timeout)
	srv := &http.Server{
		Addr:         s.cfg.Addr,
		Handler:      corsHandler,
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	go func() {
		s.log.Info("HTTP server starting", slog.String("addr", s.cfg.Addr))

		if err := srv.ListenAndServe(); err != nil {
			s.log.Error("HTTP server failed", slog.Any("error", err))
			return
		}
	}()

	return func() { _ = srv.Close() }
}

func NewRestAPI(cfg *config.HTTP, svc *service.ServiceURL, log *slog.Logger) *Server {
	return &Server{
		cfg: cfg,
		svc: svc,
		log: log,
	}
}
