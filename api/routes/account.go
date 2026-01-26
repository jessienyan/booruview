package routes

import (
	"net/http"
)

type jsonString string

// Marshal the data as-is without adding string quotes. The data is stored as a JSON
// string in the DB.
func (s jsonString) MarshalJSON() ([]byte, error) {
	return []byte(s), nil
}

type AccountResponse struct {
	Data     jsonString `json:"data"`
	Username string     `json:"username"`
}

func AccountHandler(w http.ResponseWriter, req *http.Request) {
	user := getUser(req)
	if user == nil {
		respondWithUnauthorized(w)
		return
	}

	respondJson(w, 200, AccountResponse{
		Data:     jsonString(user.Data.Data),
		Username: user.User.Username,
	})
}
