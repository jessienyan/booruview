package routes

import (
	"net/http"

	"github.com/pkg/errors"
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
					err = errors.Errorf("%v", recoverErr)
				}
				respondWithInternalError(w, errors.Wrap(err, "recovered from panic"))
			}
		}()

		next.ServeHTTP(w, req)
	})
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	bodyLen int
}

func (w *wrappedResponseWriter) Write(data []byte) (int, error) {
	w.bodyLen = len(data)
	return w.ResponseWriter.Write(data)
}

// Checks for an empty response body and logs it
func EmptyResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wrapped := &wrappedResponseWriter{ResponseWriter: w}
		next.ServeHTTP(wrapped, req)

		if wrapped.bodyLen == 0 {
			log.Warn().Str("url", req.URL.String()).Msg("empty response body")
		}
	})
}
