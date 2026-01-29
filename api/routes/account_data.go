package routes

import (
	"encoding/json"
	"io"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type AccountResponse struct {
	Data     models.UserDataJSON `json:"data"`
	Username string              `json:"username"`
}

func AccountDataHandler(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		respondWithUnauthorized(w)
		return
	}

	data, err := user.Data.ParseJSON()
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	if req.Method == "PATCH" {
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

		var form models.UserDataJSON
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

		if changed {
			newData, err := data.MarshalJSON()
			if err != nil {
				respondWithInternalError(w, err)
				return
			}

			db := models.New(api.UserDB())
			params := models.UpdateUserDataParams{
				Data:   string(newData),
				UserID: user.User.ID,
			}
			if err := db.UpdateUserData(req.Context(), params); err != nil {
				respondWithInternalError(w, err)
				return
			}

			log.Info().Any("data", data).Int64("userid", user.User.ID).Msg("updated user data")

			respondJson(w, 200, data)
			return
		}
	}

	respondJson(w, 200, AccountResponse{
		Data:     data,
		Username: user.User.Username,
	})
}
