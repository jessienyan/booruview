package api

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	AppVersion = "unset" // embedded using flags at build time, check Dockerfile
	ValkeyAddr = os.Getenv("VALKEY_ADDR")
	SecretKey  = []byte(os.Getenv("SECRET_KEY"))

	// Optional
	GelbooruUserIds = []string(nil)
	GelbooruApiKeys = []string(nil)
)

func init() {
	ok := true

	if len(SecretKey) == 0 {
		log.Error().Msg("SECRET_KEY cannot be blank")
		ok = false
	}

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
		log.Error().Msg("number of gelbooru userids and apikeys does not match")
		ok = false
	}

	if !ok {
		os.Exit(1)
	}
}
