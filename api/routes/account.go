package routes

import (
	"encoding/json"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/rs/zerolog/log"
)

type accountDeleteParams struct {
	// This is a simple failsafe against accidentally sending a DELETE request
	Confirm bool `json:"permanently_delete_account"`
}

func AccountDeleteHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, accountDeleteCost) {
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
	defer req.Body.Close()

	var form accountDeleteParams
	if err := json.Unmarshal(body, &form); err != nil {
		log.Err(err).Msg("failed to parse form")
		respondWithBadRequest(w, "failed to parse json: "+err.Error())
		return
	}

	if !form.Confirm {
		log.Warn().Msg("possible accidental account deletion caught")
		respondWithBadRequest(w, "include a JSON body with permanently_delete_account = true")
		return
	}

	db := api.UserDB()
	tx, err := db.Begin()
	if err != nil {
		respondWithInternalError(w, err)
		return
	}
	defer tx.Rollback()

	query := models.New(db).WithTx(tx)
	user := GetUser(req)
	if err := query.DeleteUserData(req.Context(), user.User.ID); err != nil {
		respondWithInternalError(w, err)
		return
	}

	if err := query.DeleteUser(req.Context(), user.User.ID); err != nil {
		respondWithInternalError(w, err)
		return
	}

	if err := tx.Commit(); err != nil {
		respondWithInternalError(w, err)
		return
	}

	// On successful DELETE, return 204 no response
	w.WriteHeader(204)
}
