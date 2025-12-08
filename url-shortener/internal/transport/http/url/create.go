package url

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"shortener/internal/dto"
	"shortener/internal/service"
	errorhandle "shortener/internal/transport/http/error"
)

func CreateURL(svc *service.ServiceURL, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body dto.CreateCodeURLRequest

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errorhandle.SendAppError(w, service.WrapErrBadRequest(err))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		code, err := svc.CreateShortURL(ctx, body.URL)
		if err != nil {
			errorhandle.SendAppError(w, err)
			return
		}

		response := dto.CreateCodeURLResponse{
			Code: code,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		_ = json.NewEncoder(w).Encode(response)
	}
}
