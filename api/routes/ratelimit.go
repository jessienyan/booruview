package routes

import (
	"net/http"
	"strconv"
	"time"

	api "codeberg.org/jessienyan/booruview"
)

const (
	postApiCostIfCacheHit  = 1
	postApiCostIfCacheMiss = 10
	settingExportApiCost   = 10
	settingImportApiCost   = 1
	tagSearchApiCost       = 1
	tagApiCost             = 2
)

func isRateLimited(w http.ResponseWriter, req *http.Request, cost int) (abort bool) {
	cb, err := api.IsRateLimited(clientIP(req), cost)
	if err != nil {
		respondWithInternalError(w, err)
		abort = true
	} else if cb.Banned() {
		banDuration := time.Until(cb.BannedUntil).Round(time.Second)
		w.Header().Add("Retry-After", strconv.Itoa(int(banDuration.Seconds())))
		respondWithRateLimited(w, banDuration)
		abort = true
	}

	return
}
