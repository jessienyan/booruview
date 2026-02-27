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
func loadUserIntoContext(w http.ResponseWriter, req *http.Request, parsedToken api.ParsedAuthToken) (newReq *http.Request, ok bool) {
	db := models.New(api.UserDB())
	user, err := db.GetUserByID(req.Context(), parsedToken.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int64("userid", parsedToken.UserID).Msg("user tried to login but account doesn't exist, probably deleted")
			respondWithUnauthorized(w)
			return nil, false
		}
		respondWithInternalError(w, err)
		return nil, false
	}

	// Reject the token if the user changed their password since the token was created.
	if user.PasswordChangedAt.Valid && user.PasswordChangedAt.Time.After(parsedToken.CreatedAt) {
		log.Info().Int64("userid", parsedToken.UserID).Time("issued", parsedToken.CreatedAt).Time("password_changed_at", user.PasswordChangedAt.Time).Msg("rejected token created before last password change")
		respondWithUnauthorized(w)
		return nil, false
	}

	data, err := db.GetUserData(req.Context(), parsedToken.UserID)
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
// If the token is invalid or missing, returns 401.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			respondWithUnauthorized(w)
			return
		}

		if parsedToken, err := api.ParseAuthToken(token); err == nil {
			var ok bool
			if req, ok = loadUserIntoContext(w, req, parsedToken); !ok {
				// loadUserIntoContext sent a response already, nothing to do here
				return
			}
		} else {
			respondWithUnauthorized(w)
			return
		}

		next.ServeHTTP(w, req)
	})
}
