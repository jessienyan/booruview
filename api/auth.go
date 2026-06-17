package api

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/argon2"
)

const (
	hashTime      = 3
	hashMemory    = 16 * 1024 // 16MB
	hashThreads   = 1
	hashKeyLength = 32
	saltLength    = 16

	sessionKeyLength = 16

	// How long users stay logged in for
	SessionTTL = time.Hour * 24 * 90

	AuthCookieName = "booruviewauth"
)

var (
	// For easy stubbing in tests
	// source: https://ekm.id.au/posts/golang-mock-time/
	Now = time.Now
)

func HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, hashTime, hashMemory, hashThreads, hashKeyLength)
}

func GenerateSalt() []byte {
	buf := make([]byte, saltLength)
	rand.Read(buf)
	return buf
}

var (
	SessionInvalid = errors.New("session is not valid")
	SessionExpired = errors.New("session has expired")
)

func GenerateSessionKey() string {
	buf := make([]byte, sessionKeyLength)
	rand.Read(buf)
	return hex.EncodeToString(buf)
}
