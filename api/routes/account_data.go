package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

const (
	maxDataSize = 4 * 1024 * 1024 // MB
)

type AccountDataResponse struct {
	models.UserDataJSON
}

func AccountDataGetHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, accountDataGetCost) {
		return
	}

	user := GetUser(req)
	data, err := user.Data.ParseJSON()
	if err != nil {
		err = errors.Wrap(err, "failed to parse JSON")
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

	// TODO: wrap this in a transaction

	user := GetUser(req)
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

type AddAccountData struct {
	FavoritePosts api.PostList             `json:"favorite_posts"`
	FavoriteTags  api.TagList              `json:"favorite_tags"`
	Blacklist     api.TagList              `json:"blacklist"`
	SearchHistory models.SearchHistoryList `json:"search_history" validate:"dive"`
}

type RemoveAccountData struct {
	FavoritePosts []int                `json:"favorite_post_ids"`
	FavoriteTags  []string             `json:"favorite_tag_names"`
	Blacklist     []string             `json:"blacklist_names"`
	SearchHistory []models.SearchQuery `json:"search_history" validate:"dive"`
}

type AccountDataPatchParams struct {
	Add    *AddAccountData    `json:"add"`
	Remove *RemoveAccountData `json:"remove"`
}

func AccountDataPatchHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, accountDataPatchCost) {
		return
	}

	// TODO: wrap this in a transaction

	user := GetUser(req)
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

	var form AccountDataPatchParams
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

	if form.Add != nil {
		if len(form.Add.Blacklist) > 0 {
			data.Blacklist = append(data.Blacklist, form.Add.Blacklist...)
			changed = true
		}

		if len(form.Add.FavoritePosts) > 0 {
			// NOTE: for compatibility, new posts are added to the beginning of the list
			data.FavoritePosts = append(form.Add.FavoritePosts, data.FavoritePosts...)
			changed = true
		}

		if len(form.Add.FavoriteTags) > 0 {
			data.FavoriteTags = append(data.FavoriteTags, form.Add.FavoriteTags...)
			changed = true
		}

		if len(form.Add.SearchHistory) > 0 {
			data.SearchHistory = append(data.SearchHistory, form.Add.SearchHistory...)
			changed = true
		}
	}

	if form.Remove != nil {
		if len(form.Remove.Blacklist) > 0 {
			data.Blacklist.Remove(form.Remove.Blacklist)
			changed = true
		}

		if len(form.Remove.FavoritePosts) > 0 {
			data.FavoritePosts.Remove(form.Remove.FavoritePosts)
			changed = true
		}

		if len(form.Remove.FavoriteTags) > 0 {
			data.FavoriteTags.Remove(form.Remove.FavoriteTags)
			changed = true
		}

		if len(form.Remove.SearchHistory) > 0 {
			data.SearchHistory.Remove(form.Remove.SearchHistory)
			changed = true
		}
	}

	if changed {
		data.Clean()

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
