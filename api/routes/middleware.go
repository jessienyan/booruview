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

type validatedSession struct {
	UserID int64
}

func validateSession(ctx context.Context, key string) (validatedSession, error) {
	if key == "" {
		return validatedSession{}, api.SessionInvalid
	}

	db := models.New(api.UserDB())
	db.DeleteExpiredSessions(ctx)

	session, err := db.GetSessionByKey(ctx, key)
	if err != nil {
		if err == sql.ErrNoRows {
			return validatedSession{}, api.SessionInvalid
		}
		return validatedSession{}, api.SessionInvalid
	}

	if session.ExpiresAt.Before(api.Now()) {
		return validatedSession{}, api.SessionExpired
	}

	return validatedSession{UserID: session.UserID}, nil
}

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
func loadUserIntoContext(w http.ResponseWriter, req *http.Request, session validatedSession) (newReq *http.Request, ok bool) {
	db := models.New(api.UserDB())
	user, err := db.GetUserByID(req.Context(), session.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Info().Int64("userid", session.UserID).Msg("user tried to login but account doesn't exist, probably deleted")
			respondWithUnauthorized(w)
			return nil, false
		}
		err = errors.Wrap(err, "failed to lookup user")
		respondWithInternalError(w, err)
		return nil, false
	}

	data, err := db.GetUserData(req.Context(), session.UserID)
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

// Adds a user to the request context if a valid session cookie was sent.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cookie, err := req.Cookie(api.AuthCookieName)
		if err != nil || cookie == nil {
			next.ServeHTTP(w, req)
			return
		}

		if session, err := validateSession(req.Context(), cookie.Value); err == nil {
			// Make sure we don't clobber `req`
			var ok bool
			if req, ok = loadUserIntoContext(w, req, session); !ok {
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
