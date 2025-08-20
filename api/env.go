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
	GelbooruUserIds     = []string(nil)
	GelbooruApiKeys     = []string(nil)
	NaughtyFingerprints = make(map[string]bool)
	FakePostHashes      = []string(nil)
)

func init() {
	if userIds := os.Getenv("GELBOORU_USERID"); userIds != "" {
		GelbooruUserIds = strings.Split(userIds, ",")
	} else {
		log.Warn().Msg("GELBOORU_USERID is not set (may be subject to rate limiting)")
	}

	if apiKeys := os.Getenv("GELBOORU_APIKEY"); apiKeys != "" {
		GelbooruApiKeys = strings.Split(apiKeys, ",")
	} else {
		log.Warn().Msg("GELBOORU_APIKEY is not set (may be subject to rate limiting)")
	}

	if len(GelbooruUserIds) != len(GelbooruApiKeys) {
		log.Fatal().Msg("number of gelbooru userids and apikeys does not match")
	}

	if fingerprints := os.Getenv("NAUGHTY_JA4H_FINGERPRINTS"); fingerprints != "" {
		for fp := range strings.SplitSeq(fingerprints, ",") {
			NaughtyFingerprints[fp] = true
		}

		log.Info().Msgf("loaded %d ja4h fingerprints", len(NaughtyFingerprints))
	}

	if hashes := os.Getenv("FAKEDATA_POST_HASHES"); hashes != "" {
		FakePostHashes = strings.Split(hashes, ",")
	}
}
