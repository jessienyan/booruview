package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type errResponse struct {
	Error string `json:"error"`
	Extra string `json:"extra,omitempty"`
}

// respondJson writes a JSON response and returns the response body as []byte.
// Panics if there is an error
func respondJson(w http.ResponseWriter, code int, data any) []byte {
	var resp []byte

	if dataAsBytes, ok := data.([]byte); ok {
		resp = dataAsBytes
	} else {
		var err error
		resp, err = json.Marshal(data)
		if err != nil {
			panic(err)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(resp); err != nil {
		panic(err)
	}

	return resp
}

func respondWithError(w http.ResponseWriter, code int, msg string, extra string) {
	respondJson(w, code, errResponse{Error: msg, Extra: extra})
}

func respondWithBadRequest(w http.ResponseWriter, msg string) {
	respondWithError(w, http.StatusBadRequest, msg, "")
}

func respondWithNotFound(w http.ResponseWriter, msg string) {
	respondWithError(w, http.StatusNotFound, msg, "")
}

func respondWithRateLimited(w http.ResponseWriter, banDuration time.Duration) {
	respondWithError(
		w,
		http.StatusTooManyRequests,
		fmt.Sprintf("Rate limited, wait %s and try again", banDuration),
		"Please do not use booruview's API for scraping. Use Gelbooru's API directly. Any form of abuse will result in a permanent IP ban.",
	)
}

func respondWithGelbooruUnavailable(w http.ResponseWriter) {
	respondWithError(w, http.StatusServiceUnavailable, "Gelbooru is currently unavailable", "")
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func respondWithInternalError(w http.ResponseWriter, err error) {
	log.Err(err).Msg("api error")
	if err, ok := err.(stackTracer); ok {
		log.Info().Msgf("Stacktrace%+v", err.StackTrace())
	}
	respondWithError(w, http.StatusInternalServerError, "An unexpected error occurred :(", "")
}
