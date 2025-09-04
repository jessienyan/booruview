package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"github.com/rs/zerolog/log"
)

type errResponse struct {
	Error string `json:"error"`
	Extra string `json:"extra,omitempty"`
}

// respondJson writes a JSON response and returns the response body as []byte
func respondJson(w http.ResponseWriter, code int, data any) []byte {
	var resp []byte

	if dataAsBytes, ok := data.([]byte); ok {
		resp = dataAsBytes
	} else {
		resp, _ = json.Marshal(data)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)

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

func respondWithInternalError(w http.ResponseWriter, err error) {
	log.Err(err).Stack().Msg("api error")
	api.LogStackTrace()
	respondWithError(w, http.StatusInternalServerError, "An unexpected error occurred :(", "")
}
