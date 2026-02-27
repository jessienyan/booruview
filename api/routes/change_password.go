package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
)

type ChangePasswordParams struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ChangePasswordResponse struct {
	AuthToken string `json:"auth_token"`
}

func ChangePasswordHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, resetPasswordCost) {
		return
	}

	user := getUser(req)
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

	// Generate a new token for the user so they aren't immediately logged out
	token, err := api.NewAuthToken(int(user.User.ID), api.AuthTokenTTL)
	if err != nil {
		err = fmt.Errorf("failed to create new auth token: %w", err)
		respondWithInternalError(w, err)
		return
	}

	respondJson(w, 200, ChangePasswordResponse{AuthToken: token})
}
