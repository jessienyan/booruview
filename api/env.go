package api

import (
	"os"
)

var (
	AppVersion = "unset" // embedded using flags at build time, check Dockerfile
	ValkeyAddr = os.Getenv("VALKEY_ADDR")

	// Optional
	GelbooruUserId = os.Getenv("GELBOORU_USERID")
	GelbooruApiKey = os.Getenv("GELBOORU_APIKEY")
)
