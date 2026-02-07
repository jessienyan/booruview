package routes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func clientIP(req *http.Request) string {
	val := req.Header.Get("X-Forwarded-For")
	if val == "" {
		log.Warn().Msg("request context is missing IP")
	}

	return val
}

// Gracefully recover from a panic
func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				var err error
				if asErr, ok := recoverErr.(error); ok {
					err = asErr
				} else {
					err = errors.Errorf("%v", recoverErr)
				}
				respondWithInternalError(w, errors.Wrap(err, "recovered from panic"))
			}
		}()

		next.ServeHTTP(w, req)
	})
}

type contextKey struct{}

var userKey contextKey

type userContextValue struct {
	User models.Users
	Data models.UserData
}

func getUser(req *http.Request) *userContextValue {
	user := req.Context().Value(userKey)
	if user == nil {
		return nil
	}
	return user.(*userContextValue)
}

// If the auth token is valid, fetches the user and adds it to the request context.
// If the token isn't valid, it does nothing.
func LoadUserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		userIDFromToken := func() (int64, bool) {
			token := req.Header.Get("Authorization")
			token = strings.TrimPrefix(token, "Bearer ")
			if token == "" {
				return 0, false
			}

			userID, err := api.ParseAuthToken(token)
			if err != nil {
				return 0, false
			}

			return int64(userID), true
		}

		loadUserIntoContext := func(userID int64) {
			db := models.New(api.UserDB())
			user, err := db.GetUserByID(req.Context(), userID)
			if err != nil {
				if err != sql.ErrNoRows {
					respondWithInternalError(w, err)
				}
				return
			}

			data, err := db.GetUserData(req.Context(), userID)
			if err == sql.ErrNoRows {
				var userData models.UserData
				userData.Set(models.UserDataJSON{})

				// Create user data if it's missing
				data, err = db.CreateUserData(req.Context(), models.CreateUserDataParams{
					Data:   userData.Data,
					UserID: user.ID,
				})
				if err != nil {
					respondWithInternalError(w, err)
					return
				}

				log.Info().Int64("userid", user.ID).Msg("created missing user data")
			} else if err != nil {
				return
			}

			req = req.WithContext(context.WithValue(req.Context(), userKey, &userContextValue{
				User: user,
				Data: data,
			}))
		}

		if id, ok := userIDFromToken(); ok {
			loadUserIntoContext(id)
		}

		next.ServeHTTP(w, req)
	})
}
