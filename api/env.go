package api

import (
	"log"
	"os"
)

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("required env '%s' is missing or empty", key)
	}
	return val
}

var (
	AppVersion = mustGetEnv("COMMIT_HASH")
	ValkeyAddr = mustGetEnv("VALKEY_ADDR")

	// Optional
	GelbooruUserId = os.Getenv("GELBOORU_USERID")
	GelbooruApiKey = os.Getenv("GELBOORU_APIKEY")
)
