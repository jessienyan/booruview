package routes

import (
	"context"
	"database/sql"
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

// If ok, fetches a user and adds them to the request context if they exist.
// Otherwise, an error response is written.
func loadUserIntoContext(w http.ResponseWriter, req *http.Request, userID int64) (newReq *http.Request, ok bool) {
	db := models.New(api.UserDB())
	user, err := db.GetUserByID(req.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int64("userid", userID).Msg("user tried to login but account doesn't exist, probably deleted")
			respondWithUnauthorized(w)
			return nil, false
		}
		respondWithInternalError(w, err)
		return nil, false
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
			return nil, false
		}

		log.Info().Int64("userid", user.ID).Msg("created missing user data")
	} else if err != nil {
		return nil, false
	}

	req = req.WithContext(context.WithValue(req.Context(), userKey, &userContextValue{
		User: user,
		Data: data,
	}))

	return req, true
}

// Adds a user to the request context if a valid auth token was sent.
// If the token is invalid, returns 401.
// If there's no token, it skips the auth check.
func MaybeAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		if token != "" {
			if userID, err := api.ParseAuthToken(token); err == nil {
				var ok bool
				if req, ok = loadUserIntoContext(w, req, userID); !ok {
					return
				}
			} else {
				respondWithUnauthorized(w)
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}
