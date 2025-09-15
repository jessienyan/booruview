package routes

import (
	"net/http"
	"strconv"
	"time"

	api "codeberg.org/jessienyan/booruview"
)

const (
	postApiCostIfCacheHit  = 1
	postApiCostIfCacheMiss = 12
	settingExportApiCost   = 10
	settingImportApiCost   = 1
	tagSearchApiCost       = 1
	tagApiCost             = 2
)

func isRateLimited(w http.ResponseWriter, req *http.Request, cost int) bool {
	cb, err := api.IsRateLimited(clientIP(req), cost)

	if err != nil {
		respondWithInternalError(w, err)
		return true
	} else if cb.Banned() {
		banDuration := time.Until(cb.BannedUntil).Round(time.Second)
		w.Header().Add("Retry-After", strconv.Itoa(int(banDuration.Seconds())))
		respondWithRateLimited(w, banDuration)
		return true
	}

	return false
}
