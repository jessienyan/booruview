package api

import (
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	AppVersion = "unset" // embedded using flags at build time, check Dockerfile
	ValkeyAddr = os.Getenv("VALKEY_ADDR")

	// Optional
	GelbooruUserIds = []string(nil)
	GelbooruApiKeys = []string(nil)
	MediaProxyHost  string
)

var (
	reProxy = regexp.MustCompile(`^https?:\/\/.+[^\/]$`)
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

	MediaProxyHost = os.Getenv("MEDIA_PROXY_HOST")
	if MediaProxyHost != "" && !reProxy.MatchString(MediaProxyHost) {
		log.Fatal().Msg("MEDIA_PROXY_HOST must either be blank, or a http/https origin, e.g. 'https://example.com'")
	}

	if len(GelbooruUserIds) != len(GelbooruApiKeys) {
		log.Fatal().Msg("number of gelbooru userids and apikeys does not match")
	}
}
