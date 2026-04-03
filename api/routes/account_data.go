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
	maxDataSize = 4 * 1024 * 1024 // 4 MB
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
	if form.SavedSearches != nil {
		changed = true
		data.SavedSearches = form.SavedSearches
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

type AddAccountData struct {
	FavoritePosts api.PostList             `json:"favorite_posts" validate:"dive"`
	FavoriteTags  api.TagList              `json:"favorite_tags" validate:"dive"`
	Blacklist     api.TagList              `json:"blacklist" validate:"dive"`
	SearchHistory models.SearchHistoryList `json:"search_history" validate:"dive"`
	SavedSearches models.SearchHistoryList `json:"saved_searches" validate:"dive"`
}

type RemoveAccountData struct {
	FavoritePostIDs  []int                     `json:"favorite_post_ids"`
	FavoriteTagNames []string                  `json:"favorite_tag_names"`
	BlacklistNames   []string                  `json:"blacklist_names"`
	SearchQueries    []models.SearchQueryNames `json:"search_queries" validate:"dive"`
	SavedQueries     []models.SearchQueryNames `json:"saved_queries" validate:"dive"`
}

type AccountDataPatchParams struct {
	Add    AddAccountData    `json:"add"`
	Remove RemoveAccountData `json:"remove"`
}

type AccountDataPatchResponse struct {
	FavoritePosts api.PostList
	FavoriteTags  api.TagList
	Blacklist     api.TagList
	SearchHistory models.SearchHistoryList
	SavedSearches models.SearchHistoryList
}

// Overrides the default marshal behavior to only include fields that are non-nil.
// This fixes an issue where a PATCH removes all the entries from a field, but the
// response is empty `{}` (due to using `omitempty`). This behavior is bug-prone if
// the caller always expects modified fields to be included in the response
func (resp AccountDataPatchResponse) MarshalJSON() ([]byte, error) {
	respMap := make(map[string]any)

	if resp.FavoritePosts != nil {
		respMap["favorite_posts"] = resp.FavoritePosts
	}
	if resp.FavoriteTags != nil {
		respMap["favorite_tags"] = resp.FavoriteTags
	}
	if resp.Blacklist != nil {
		respMap["blacklist"] = resp.Blacklist
	}
	if resp.SearchHistory != nil {
		respMap["search_history"] = resp.SearchHistory
	}
	if resp.SavedSearches != nil {
		respMap["saved_searches"] = resp.SavedSearches
	}

	return json.Marshal(respMap)
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

	response := AccountDataPatchResponse{}

	if len(form.Add.Blacklist) > 0 {
		data.Blacklist.Add(form.Add.Blacklist)
		response.Blacklist = data.Blacklist
	}
	if len(form.Add.FavoritePosts) > 0 {
		data.FavoritePosts.Add(form.Add.FavoritePosts)
		response.FavoritePosts = data.FavoritePosts
	}
	if len(form.Add.FavoriteTags) > 0 {
		data.FavoriteTags.Add(form.Add.FavoriteTags)
		response.FavoriteTags = data.FavoriteTags
	}
	if len(form.Add.SearchHistory) > 0 {
		data.SearchHistory.Add(form.Add.SearchHistory)
		response.SearchHistory = data.SearchHistory
	}
	if len(form.Add.SavedSearches) > 0 {
		data.SavedSearches.Add(form.Add.SavedSearches)
		response.SavedSearches = data.SavedSearches
	}

	if len(form.Remove.BlacklistNames) > 0 {
		data.Blacklist.Remove(form.Remove.BlacklistNames)
		response.Blacklist = data.Blacklist
	}
	if len(form.Remove.FavoritePostIDs) > 0 {
		data.FavoritePosts.Remove(form.Remove.FavoritePostIDs)
		response.FavoritePosts = data.FavoritePosts
	}
	if len(form.Remove.FavoriteTagNames) > 0 {
		data.FavoriteTags.Remove(form.Remove.FavoriteTagNames)
		response.FavoriteTags = data.FavoriteTags
	}
	if len(form.Remove.SearchQueries) > 0 {
		queries := make([]string, 0, len(form.Remove.SearchQueries))
		for _, q := range form.Remove.SearchQueries {
			q.Clean()
			queries = append(queries, q.Tags())
		}
		data.SearchHistory.Remove(queries)
		response.SearchHistory = data.SearchHistory
	}
	if len(form.Remove.SavedQueries) > 0 {
		queries := make([]string, 0, len(form.Remove.SavedQueries))
		for _, q := range form.Remove.SavedQueries {
			q.Clean()
			queries = append(queries, q.Tags())
		}
		data.SavedSearches.Remove(queries)
		response.SavedSearches = data.SavedSearches
	}

	if response.FavoritePosts != nil || response.FavoriteTags != nil || response.Blacklist != nil || response.SearchHistory != nil || response.SavedSearches != nil {
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

	respondJson(w, 200, response)
}
