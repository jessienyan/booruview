package api

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	hashTime      = 3
	hashMemory    = 16 * 1024 // 16MB
	hashThreads   = 1
	hashKeyLength = 32
	saltLength    = 16
)

func HashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, hashTime, hashMemory, hashThreads, hashKeyLength)
}

func GenerateSalt() []byte {
	buf := make([]byte, saltLength)
	rand.Read(buf)
	return buf
}
