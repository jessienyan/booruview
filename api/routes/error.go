package routes

import (
	"encoding/json"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"github.com/rs/zerolog/log"
)

type errResponse struct {
	Error string `json:"error"`
}

func handle400Error(w http.ResponseWriter, msg string) {
	handle4xxError(w, 400, msg)
}

func handle4xxError(w http.ResponseWriter, code int, msg string) {
	resp, _ := json.Marshal(errResponse{Error: msg})
	w.WriteHeader(code)
	w.Write(resp)
}

func handle429Error(w http.ResponseWriter) {
	handle4xxError(w, 429, "Rate limited, wait a few minutes and try again")
}

func handleGelbooruUnavailable(w http.ResponseWriter) {
	handle4xxError(w, 503, "Gelbooru is currently unavailable")
}

func handleError(w http.ResponseWriter, err error) {
	log.Err(err).Stack().Msg("api error")
	api.LogStackTrace()
	w.WriteHeader(http.StatusInternalServerError)
}
