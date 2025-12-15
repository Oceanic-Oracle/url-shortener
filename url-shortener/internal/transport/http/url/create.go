package url

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"shortener/internal/dto"
	"shortener/internal/logctx"
	"shortener/internal/service"
	httperror "shortener/internal/transport/http/error"
)

func CreateURL(svc *service.ServiceURL, host string, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.InfoContext(ctx, "Received Create URL request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		var body dto.CreateCodeURLRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			log.WarnContext(ctx, "Invalid JSON in request body", slog.Any("error", err))
			httperror.SendHTTPError(w, httperror.ErrBadRequest)

			return
		}

		ctx = logctx.WithURL(ctx, body.URL)

		if err := body.Validate(host); err != nil {
			log.WarnContext(ctx, "Invalid URL in request body", slog.Any("error", err))
			httperror.SendHTTPError(w, httperror.ErrBadRequest)

			return
		}

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		code, err := svc.CreateShortURL(ctx, body.URL)
		if err != nil {
			log.ErrorContext(ctx, "Error creating shortcode", slog.Any("error", err))
			httperror.SendHTTPError(w, err)

			return
		}

		response := dto.CreateCodeURLResponse{
			Code: code,
		}

		ctx = logctx.WithCode(ctx, code)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)

		log.InfoContext(ctx, "Create short code successfully")
	}
}
