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

				// This happens when the connection closes while we're writing a response.
				// Usually it means the user refreshed the page or closed their browser
				if strings.Contains(err.Error(), "write: broken pipe") {
					return
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

func GetUser(req *http.Request) *userContextValue {
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
		err = errors.Wrap(err, "failed to lookup user")
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
			err = errors.Wrap(err, "failed to create user data")
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
// Looks for the token in the Authorization header, then as a cookie.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var token string

		token = strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")

		if token == "" {
			cookie, _ := req.Cookie(api.AuthCookieName)
			if cookie != nil {
				token = cookie.Value
			}
		}

		if token != "" {
			if parsedToken, err := api.ParseAuthToken(token); err == nil {
				// Make sure we don't clobber `req`
				var ok bool
				if req, ok = loadUserIntoContext(w, req, parsedToken); !ok {
					// loadUserIntoContext sent a response already, nothing to do here
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

// Aborts the request with a 401 if the user wasn't authenticated
func RequireAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		u := GetUser(req)
		if u == nil {
			respondWithUnauthorized(w)
			return
		}

		next.ServeHTTP(w, req)
	})
}
