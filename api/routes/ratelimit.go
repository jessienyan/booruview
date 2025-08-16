package routes

import (
	"net/http"

	api "github.com/jessienyan/booruview"
)

func isRateLimited(w http.ResponseWriter, req *http.Request, cost int) (abort bool) {
	banned, err := api.IsRateLimited(clientIP(req), cost, req.Header.Get("Ja4h"))
	if err != nil {
		handleError(w, err)
		abort = true
	} else if banned {
		handle429Error(w)
		abort = true
	}

	return
}
