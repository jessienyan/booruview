package routes

import (
	"net/http"

	"codeberg.org/jessienyan/booruview/models"
)

type AccountResponse struct {
	Data     models.UserDataJSON `json:"data"`
	Username string              `json:"username"`
}

func AccountHandler(w http.ResponseWriter, req *http.Request) {
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

	respondJson(w, 200, AccountResponse{
		Data:     data,
		Username: user.User.Username,
	})
}
