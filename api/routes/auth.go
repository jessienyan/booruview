package routes

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	api "codeberg.org/jessienyan/booruview"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

const (
	hashTime      = 3
	hashMemory    = 16 * 1024 // 16MB
	hashThreads   = 1
	hashKeyLength = 32
	saltLength    = 16

	minUsernameLength = 3
	maxUsernameLength = 16
	minPasswordLength = 8
)

var (
	reUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

func hashPassword(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, hashTime, hashMemory, hashThreads, hashKeyLength)
}

func generateSalt() []byte {
	buf := make([]byte, saltLength)
	rand.Read(buf)
	return buf
}

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if req.Header.Get("Content-Type") != "application/json" {
		respondWithBadRequest(w, "expected Content-Type header to be application/json")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	var params CreateUserParams
	if err := json.Unmarshal(body, &params); err != nil {
		respondWithBadRequest(w, "json body is not valid")
		return
	}

	params.Username = strings.TrimSpace(params.Username)

	if !reUsername.MatchString(params.Username) {
		respondWithBadRequest(w, "username can only contain letters, numbers, hyphens, and underscores")
		return
	}

	if len(params.Username) < minUsernameLength || len(params.Username) > maxUsernameLength {
		respondWithBadRequest(w, fmt.Sprintf("username must be %d-%d characters", minUsernameLength, maxUsernameLength))
		return
	}

	if len(params.Password) < minPasswordLength {
		respondWithBadRequest(w, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
		return
	}

	db := api.UserDB()
	row := db.QueryRow(`SELECT 1 FROM users WHERE LOWER(username) = ?`, strings.ToLower(params.Username))
	err = row.Err()
	usernameTaken := err == nil
	otherError := err != nil && err != sql.ErrNoRows

	if usernameTaken {
		respondWithBadRequest(w, "username is already taken")
		return
	} else if otherError {
		respondWithInternalError(w, err)
		return
	}

}
