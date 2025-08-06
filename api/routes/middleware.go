package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	api "github.com/jessienyan/booruview"
)

func clientIP(req *http.Request) string {
	val := req.Header.Get("X-Forwarded-For")
	if val == "" {
		log.Warn().Msg("request context is missing IP")
	}

	return val
}

// Gracefully recover from a panic
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Err(fmt.Errorf("%v", err)).Msg("recovered from panic")
				api.LogStackTrace()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, req)
	})
}
