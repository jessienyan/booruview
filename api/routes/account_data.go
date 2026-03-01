package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

const (
	maxDataSize = 2 * 1024 * 1024 // MB
)

type AccountDataResponse struct {
	models.UserDataJSON
}

func AccountDataGetHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, accountDataGetCost) {
		return
	}

	user := getUser(req)
	data, err := user.Data.ParseJSON()
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	respondJson(w, 200, AccountDataResponse{data})
}

type AccountDataPutParams struct {
	models.UserDataJSON
}

func AccountDataPutHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, accountDataPutCost) {
		return
	}

	user := getUser(req)
	data, err := user.Data.ParseJSON()
	if err != nil {
		respondWithInternalError(w, err)
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

	if len(body) > maxDataSize {
		respondWithBadRequest(w, fmt.Sprintf("body is too large (max %d bytes)", maxDataSize))
		return
	}

	var form AccountDataPutParams
	if err := json.Unmarshal(body, &form); err != nil {
		log.Err(err).Msg("failed to parse form")
		respondWithBadRequest(w, "form is not valid")
		return
	}

	if err := api.V.Struct(form); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			log.Err(err).Msg("validation failed")
			respondWithBadRequest(w, validationErrs.Error())
			return
		}
		respondWithInternalError(w, err)
		return
	}

	changed := false

	if form.Blacklist != nil {
		changed = true
		data.Blacklist = form.Blacklist
	}
	if form.FavoritePosts != nil {
		changed = true
		data.FavoritePosts = form.FavoritePosts
	}
	if form.FavoriteTags != nil {
		changed = true
		data.FavoriteTags = form.FavoriteTags
	}
	if form.SearchHistory != nil {
		changed = true
		data.SearchHistory = form.SearchHistory
	}

	if changed {
		if err := user.Data.Set(data); err != nil {
			respondWithInternalError(w, err)
			return
		}

		db := models.New(api.UserDB())
		params := models.UpdateUserDataParams{
			Data:   user.Data.Data,
			UserID: user.User.ID,
		}

		if err := db.UpdateUserData(req.Context(), params); err != nil {
			respondWithInternalError(w, err)
			return
		}

		log.Info().Int64("userid", user.User.ID).Msg("updated user data")
	}

	respondJson(w, 200, AccountDataResponse{data})
}
