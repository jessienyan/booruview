package api

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	AppVersion = "unset" // embedded using flags at build time, check Dockerfile
	ValkeyAddr = os.Getenv("VALKEY_ADDR")

	// Optional
	GelbooruUserId      = os.Getenv("GELBOORU_USERID")
	GelbooruApiKey      = os.Getenv("GELBOORU_APIKEY")
	NaughtyFingerprints = make(map[string]bool)
)

func init() {
	fingerprints := os.Getenv("NAUGHTY_RATELIMIT_FINGERPRINTS")
	if fingerprints != "" {
		for fp := range strings.SplitSeq(fingerprints, ",") {
			NaughtyFingerprints[fp] = true
		}

		log.Info().Msgf("loaded %d ja4h fingerprints", len(NaughtyFingerprints))
	}
}
