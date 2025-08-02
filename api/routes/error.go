package routes

import (
	"encoding/json"
	"net/http"

	api "github.com/jessienyan/booruview"
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
	resp, _ := json.Marshal(errResponse{Error: "Rate limited, wait a few minutes and try again"})
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write(resp)
}

func handleError(w http.ResponseWriter, err error) {
	log.Err(err).Stack().Msg("api error")
	api.LogStackTrace()
	w.WriteHeader(http.StatusInternalServerError)
}
