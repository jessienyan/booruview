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
	"github.com/rs/zerolog/log"
)

const (
	minUsernameLength = 3
	maxUsernameLength = 16
	minPasswordLength = 8
)

var (
	reUsername = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AuthToken string `json:"auth_token"`
}

// Creates a new user account
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

	log.Info().Int("id", int(u.ID)).Str("username", u.Username).Msg("user registered")

	token, err := api.NewAuthToken(int(u.ID), api.AuthTokenTTL)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	respondJson(w, 200, RegisterResponse{token})
}
