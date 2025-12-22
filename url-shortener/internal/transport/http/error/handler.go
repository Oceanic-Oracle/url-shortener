package httperror

import (
	"encoding/json"
	"errors"
	"net/http"

	"shortener/internal/service"
)

var errMap = map[error]func(error) *httpError{
	service.ErrURLNotFound:  wrapErrNotFound,
	service.ErrURLCollision: wrapErrConflict,
	service.ErrStorage:      wrapErrInternalServer,
}

func SendHTTPError(w http.ResponseWriter, err error) {
	var httpErr *httpError

	for targetErr, wrapper := range errMap {
		if errors.Is(err, targetErr) {
			httpErr = wrapper(err)
			sendError(w, httpErr)

			return
		}
	}

	if errors.Is(err, ErrBadRequest) {
		httpErr = wrapErrBadRequest(err)
	} else {
		httpErr = wrapErrInternalServer(err)
	}

	sendError(w, httpErr)
}

func sendError(w http.ResponseWriter, err *httpError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	_ = json.NewEncoder(w).Encode(err)
}
