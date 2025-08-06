package routes

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"

	"github.com/rs/zerolog/log"

	api "github.com/jessienyan/booruview"
)

type contextKey string

var (
	ipKey = contextKey("ip")
)

func clientIP(req *http.Request) string {
	val, ok := req.Context().Value(ipKey).(string)
	if !ok {
		log.Warn().Msg("request context is missing IP")
		return ""
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

func IPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var ip string
		addr, err := netip.ParseAddrPort(req.RemoteAddr)

		if err != nil {
			log.Err(err).Str("remoteAddr", req.RemoteAddr).Msg("failed to parse remote address")
			ip = req.RemoteAddr
		} else {
			ip = addr.Addr().String()
		}

		req = req.WithContext(context.WithValue(req.Context(), ipKey, ip))

		next.ServeHTTP(w, req)
	})
}
