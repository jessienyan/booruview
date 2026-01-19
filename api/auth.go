package api

import (
	"crypto/rand"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/argon2"
)

const (
	hashTime      = 3
	hashMemory    = 16 * 1024 // 16MB
	hashThreads   = 1
	hashKeyLength = 32
	saltLength    = 16

	// How long users stay logged in for
	AuthTokenTTL = time.Hour * 24 * 90
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
	AuthTokenInvalid = errors.New("auth token is not valid")
	AuthTokenExpired = errors.New("auth token has expired")
)

func NewAuthToken(userID int, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.Itoa(userID),
		"iat": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(ttl)),
	})
	return token.SignedString(SecretKey)
}

// Parses the auth token to grab the user ID it belongs to. Also verifies it hasn't
// expired and that it wasn't tampered with.
func ParseAuthToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithExpirationRequired())
	if err != nil {
		if err == jwt.ErrTokenExpired {
			return 0, AuthTokenExpired
		}
		log.Err(err).Str("token", tokenString).Msg("error parsing token")
		return 0, AuthTokenInvalid
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Error().Str("token", tokenString).Msg("auth token claims are not jwt.MapClaims")
		return 0, AuthTokenInvalid
	}

	sub, err := claims.GetSubject()
	if err != nil || sub == "" {
		log.Error().Str("token", tokenString).Msg("missing subject")
		return 0, AuthTokenInvalid
	}

	uid, err := strconv.Atoi(sub)
	if err != nil {
		log.Error().Str("sub", sub).Msg("invalid subject")
		return 0, AuthTokenInvalid
	}

	return uid, nil
}
