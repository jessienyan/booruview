package routes

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
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
			if recoverErr := recover(); recoverErr != nil {
				var err error
				if asErr, ok := recoverErr.(error); ok {
					err = asErr
				} else {
					err = fmt.Errorf("%v", recoverErr)
				}
				respondWithInternalError(w, fmt.Errorf("recovered from panic: %w", err))
			}
		}()

		next.ServeHTTP(w, req)
	})
}
