package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 16
	minPasswordLength = 8
)

var (
	reUsername = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Username string `json:"username"`
}

// Creates a new user account
func RegisterHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, registerCost) {
		return
	}

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

	if !reUsername.MatchString(params.Username) {
		respondWithBadRequest(w, "username can only contain letters or numbers")
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

	db := models.New(api.UserDB())
	_, err = db.GetUser(req.Context(), params.Username)
	usernameTaken := err == nil
	usernameAvailable := err == sql.ErrNoRows
	otherError := !usernameAvailable && !usernameTaken

	if usernameTaken {
		respondWithBadRequest(w, "username is already taken")
		return
	} else if otherError {
		respondWithInternalError(w, err)
		return
	}

	salt := api.GenerateSalt()
	passHash := api.HashPassword(params.Password, salt)

	u, err := db.CreateUser(req.Context(), models.CreateUserParams{
		Username:     params.Username,
		Password:     passHash,
		PasswordSalt: salt,
	})
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	_, err = db.CreateUserData(req.Context(), models.CreateUserDataParams{Data: "", UserID: u.ID})
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	log.Info().Str("user", u.String()).Msg("user registered")

	sessionKey := api.GenerateSessionKey()
	if _, err := db.CreateSession(req.Context(), models.CreateSessionParams{
		Key:       sessionKey,
		UserID:    u.ID,
		ExpiresAt: api.Now().Add(api.SessionTTL),
	}); err != nil {
		respondWithInternalError(w, errors.Wrap(err, "failed to create session"))
		return
	}

	// TODO: cookie jar
	w.Header().Add(
		"Set-Cookie",
		fmt.Sprintf(
			"%s=%s; Max-Age=%d; Path=/; SameSite=strict; HttpOnly",
			api.AuthCookieName,
			sessionKey,
			int(api.SessionTTL.Seconds()),
		),
	)

	respondJson(w, 200, RegisterResponse{Username: u.Username})
}
