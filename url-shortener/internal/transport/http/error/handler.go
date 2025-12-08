package errorhandle

import (
	"encoding/json"
	"errors"
	"net/http"

	"shortener/internal/service"
)

func SendAppError(w http.ResponseWriter, err error) {
	var appError *service.AppError
	if !errors.As(err, &appError) {
		err = service.WrapErrInternalServer(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appError.Status)
	_ = json.NewEncoder(w).Encode(err)
}
