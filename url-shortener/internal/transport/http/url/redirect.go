package url

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi"

	"shortener/internal/logctx"
	"shortener/internal/service"
	httperror "shortener/internal/transport/http/error"
)

func RedirectURL(svc *service.ServiceURL, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		log.InfoContext(ctx, "Received Create URL request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
		)

		code := chi.URLParam(r, "code")
		if code == "" {
			httperror.SendHTTPError(w, httperror.ErrBadRequest)

			return
		}

		ctx = logctx.WithCode(ctx, code)

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		url, err := svc.GetLongURL(ctx, code)
		if err != nil {
			log.ErrorContext(ctx, "Failed to get URL", slog.Any("err", err))
			httperror.SendHTTPError(w, err)

			return
		}

		ctx = logctx.WithURL(ctx, url)

		http.Redirect(w, r, url, http.StatusFound)

		log.InfoContext(ctx, "Redirect successfully")
	}
}
