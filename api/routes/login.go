package routes

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AuthToken string `json:"auth_token"`
	Username  string `json:"username"`
}

// Login to an account and receive an auth token
func LoginHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, loginCost) {
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		respondWithBadRequest(w, "expected Content-Type header to be application/json")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		err = errors.Wrap(err, "failed to read request body")
		respondWithInternalError(w, err)
		return
	}

	var params LoginParams
	if err := json.Unmarshal(body, &params); err != nil {
		respondWithBadRequest(w, "failed to parse json: "+err.Error())
		return
	}

	if params.Username == "" || params.Password == "" {
		respondWithBadRequest(w, "username and password are required")
		return
	}

	if !reUsername.MatchString(params.Username) {
		respondWithBadRequest(w, "username can only contain letters, numbers, hyphens, and underscores")
		return
	}

	db := models.New(api.UserDB())
	u, err := db.GetUser(req.Context(), params.Username)
	doesntExist := false

	if err != nil {
		if err == sql.ErrNoRows {
			doesntExist = true
		} else {
			err = errors.Wrap(err, "failed to get user")
			respondWithInternalError(w, err)
			return
		}
	}

	if doesntExist || !bytes.Equal(api.HashPassword(params.Password, u.PasswordSalt), u.Password) {
		respondWithBadRequest(w, "Username or password is incorrect")
		return
	}

	u.LastLogin.Time = time.Now()
	u.LastLogin.Valid = true

	err = db.UserLoggedIn(req.Context(), models.UserLoggedInParams{
		LastLogin: u.LastLogin,
		ID:        u.ID,
	})
	if err != nil {
		err = errors.Wrap(err, "failed to update user")
		respondWithInternalError(w, err)
		return
	}

	token, err := api.NewAuthToken(int(u.ID), api.AuthTokenTTL)
	if err != nil {
		err = errors.Wrap(err, "failed to create auth token")
		respondWithInternalError(w, err)
		return
	}

	w.Header().Add("Set-Cookie", fmt.Sprintf("%s=%s; Max-Age=%d; SameSite=strict; HttpOnly", api.AuthCookieName, token, int(api.AuthTokenTTL.Seconds())))
	log.Info().Str("user", u.String()).Msg("user logged in")
	respondJson(w, 200, LoginResponse{AuthToken: token, Username: u.Username})
}
