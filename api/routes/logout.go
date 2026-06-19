package routes

import (
	"fmt"
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/rs/zerolog/log"
)

func LogoutHandler(w http.ResponseWriter, req *http.Request) {
	user := GetUser(req)
	if user == nil {
		respondWithUnauthorized(w)
		return
	}

	cookie, err := req.Cookie(api.AuthCookieName)
	if err != nil {
		respondWithUnauthorized(w)
		return
	}

	db := models.New(api.UserDB())
	if err := db.DeleteSessionByKey(req.Context(), cookie.Value); err != nil {
		err = fmt.Errorf("failed to delete session: %w", err)
		respondWithInternalError(w, err)
		return
	}

	api.RemoveAuthCookie(w)
	log.Info().Str("user", user.User.String()).Msg("user logged out")
	w.WriteHeader(http.StatusNoContent)
}
