package routes

import (
	"fmt"
	"net/http"
	"net/netip"

	"github.com/rs/zerolog/log"

	api "github.com/jessienyan/booruview"
)

// Gracefully recover from a panic
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Err(fmt.Errorf("%v", err)).Msg("recovered from panic")
				api.LogStackTrace()
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string
		addr, err := netip.ParseAddrPort(r.RemoteAddr)

		if err != nil {
			log.Err(err).Str("remoteAddr", r.RemoteAddr).Msg("failed to parse remote address")
			ip = r.RemoteAddr
		} else {
			ip = addr.Addr().String()
		}

		if api.IsRateLimited(ip) {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
