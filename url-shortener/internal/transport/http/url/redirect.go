package url

import (
	"context"
	"log/slog"
	"net/http"
	"shortener/internal/service"
	errorhandle "shortener/internal/transport/http/error"
	"time"

	"github.com/go-chi/chi"
)

func RedirectURL(svc *service.ServiceURL, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := chi.URLParam(r, "code")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		url, err := svc.GetLongURL(ctx, code)
		if err != nil {
			errorhandle.SendAppError(w, err)
			return
		}

		http.Redirect(w, r, url, http.StatusFound)
	}
}