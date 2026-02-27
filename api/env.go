package api

import (
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
)

var (
	AppVersion   = "unset" // embedded using flags at build time, check Dockerfile
	ValkeyAddr   = os.Getenv("VALKEY_ADDR")
	SecretKey    = []byte(os.Getenv("SECRET_KEY"))
	DatabasePath = os.Getenv("DATABASE_PATH")

	// Optional
	GelbooruUserIds = []string(nil)
	GelbooruApiKeys = []string(nil)

	UseMediaProxy  bool
	MediaProxyHost string
)

var (
	reProxy = regexp.MustCompile(`^https?:\/\/.+[^\/]$`)
)

func LoadDatabaseEnv() {
	if DatabasePath == "" {
		DatabasePath = "database/sqlite.db"
		log.Warn().Msgf("env DATABASE_PATH is not set, defaulting to %s", DatabasePath)
	}
}

func LoadEnv() {
	ok := true

	LoadDatabaseEnv()

	if len(SecretKey) == 0 {
		log.Error().Msg("env SECRET_KEY cannot be blank")
		ok = false
	}

	if userIds := os.Getenv("GELBOORU_USERID"); userIds != "" {
		GelbooruUserIds = strings.Split(userIds, ",")
	} else {
		log.Warn().Msg("env GELBOORU_USERID is not set (may be subject to rate limiting)")
	}

	if apiKeys := os.Getenv("GELBOORU_APIKEY"); apiKeys != "" {
		GelbooruApiKeys = strings.Split(apiKeys, ",")
	} else {
		log.Warn().Msg("env GELBOORU_APIKEY is not set (may be subject to rate limiting)")
	}

	UseMediaProxy = os.Getenv("USE_MEDIA_PROXY") == "1"
	if UseMediaProxy {
		MediaProxyHost = os.Getenv("MEDIA_PROXY_HOST")
		if MediaProxyHost != "" && !reProxy.MatchString(MediaProxyHost) {
			log.Fatal().Msg("env MEDIA_PROXY_HOST must either be blank, or a http/https origin, e.g. 'https://example.com'")
		}
	}

	if len(GelbooruUserIds) != len(GelbooruApiKeys) {
		log.Error().Msg("number of gelbooru userids and apikeys does not match")
		ok = false
	}

	if !ok {
		os.Exit(1)
	}
}
