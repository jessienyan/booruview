package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/rs/zerolog/log"
)

type ChangePasswordParams struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

func ChangePasswordHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, resetPasswordCost) {
		return
	}

	user := GetUser(req)
	if user == nil {
		respondWithUnauthorized(w)
		return
	}

	if req.Header.Get("Content-Type") != "application/json" {
		respondWithBadRequest(w, "expected content-type to be application/json")
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	var params ChangePasswordParams
	if err := json.Unmarshal(body, &params); err != nil {
		respondWithBadRequest(w, "json body is not valid")
		return
	}

	if !bytes.Equal(api.HashPassword(params.CurrentPassword, user.User.PasswordSalt), user.User.Password) {
		respondWithBadRequest(w, "password is incorrect")
		return
	}

	if len(params.NewPassword) < minPasswordLength {
		respondWithBadRequest(w, fmt.Sprintf("password must be at least %d characters", minPasswordLength))
		return
	}

	user.User.Password = api.HashPassword(params.NewPassword, user.User.PasswordSalt)
	db := models.New(api.UserDB())
	err = db.UpdateUserPassword(req.Context(), models.UpdateUserPasswordParams{
		Password: user.User.Password,
		ID:       user.User.ID,
	})

	if err != nil {
		err = fmt.Errorf("failed to update user password: %w", err)
		respondWithInternalError(w, err)
		return
	}

	log.Info().Str("user", user.User.String()).Msg("user changed password")

	// Invalidate all existing sessions and create a new one
	if err := db.DeleteUserSessions(req.Context(), user.User.ID); err != nil {
		err = fmt.Errorf("failed to delete user sessions: %w", err)
		respondWithInternalError(w, err)
		return
	}

	sessionKey := api.GenerateSessionKey()
	_, err = db.CreateSession(req.Context(), models.CreateSessionParams{
		Key:       sessionKey,
		UserID:    user.User.ID,
		ExpiresAt: api.Now().Add(api.SessionTTL),
	})
	if err != nil {
		err = fmt.Errorf("failed to create session: %w", err)
		respondWithInternalError(w, err)
		return
	}

	api.SetAuthCookie(w, sessionKey)
	respondJson(w, http.StatusOK, map[string]string{})
}
