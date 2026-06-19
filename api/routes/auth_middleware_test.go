package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"codeberg.org/jessienyan/booruview/routes"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

var (
	authMiddlewareUser     models.Users
	authMiddlewareUsername = "authmiddlewaretest"
	authMiddlewarePassword = "pass123"
)

func init() {
	testutil.Setup()
	api.InitUserDatabase()
	authMiddlewareUser, _ = testutil.CreateUser(authMiddlewareUsername, authMiddlewarePassword)
}

func sessionCookie(key string) *http.Cookie {
	return &http.Cookie{
		Name:  api.AuthCookieName,
		Value: key,
	}
}

func TestAuthMiddleware_Ok(t *testing.T) {
	testutil.Flush()

	sessionKey := testutil.CreateSession(authMiddlewareUser.ID)

	req := httptest.NewRequest("GET", "/faked", nil)
	req.AddCookie(sessionCookie(sessionKey))
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		u := routes.GetUser(r).User
		require.Equal(t, authMiddlewareUser, u)
	}))
	handler.ServeHTTP(rec, req)

	require.True(t, called)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("GET", "/faked", nil)
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		u := routes.GetUser(r)
		require.Nil(t, u)
	}))
	handler.ServeHTTP(rec, req)

	require.True(t, called)
}

func TestAuthMiddleware_InvalidSession(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("GET", "/faked", nil)
	req.AddCookie(sessionCookie("invalidsession"))
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// invalid sessions are silently ignored when auth is optional
		called = true
	}))
	handler.ServeHTTP(rec, req)

	require.True(t, called)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_ExpiredSession(t *testing.T) {
	testutil.Flush()

	// Create a session that's already expired
	db := models.New(api.UserDB())
	key := api.GenerateSessionKey()
	db.CreateSession(t.Context(), models.CreateSessionParams{
		Key:       key,
		UserID:    authMiddlewareUser.ID,
		ExpiresAt: api.Now().Add(-time.Hour),
	})

	req := httptest.NewRequest("GET", "/faked", nil)
	req.AddCookie(sessionCookie(key))
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// invalid sessions are silently ignored when auth is optional
		called = true
	}))
	handler.ServeHTTP(rec, req)

	require.Empty(t, rec.Result().Cookies()) // should clear cookie
	require.True(t, called)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestAuthMiddleware_UserDeleted(t *testing.T) {
	testutil.Flush()

	// Create a session for a non-existent user
	db := models.New(api.UserDB())
	key := api.GenerateSessionKey()
	db.CreateSession(t.Context(), models.CreateSessionParams{
		Key:       key,
		UserID:    123456789,
		ExpiresAt: api.Now().Add(time.Minute),
	})

	req := httptest.NewRequest("GET", "/faked", nil)
	req.AddCookie(sessionCookie(key))
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestRequireAuthMiddleware_Ok(t *testing.T) {
	testutil.Flush()

	sessionKey := testutil.CreateSession(authMiddlewareUser.ID)

	req := httptest.NewRequest("GET", "/faked", nil)
	req.AddCookie(sessionCookie(sessionKey))
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(routes.RequireAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	})))
	handler.ServeHTTP(rec, req)

	require.True(t, called)
}

func TestRequireAuthMiddleware_NoAuth(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("GET", "/faked", nil)
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(routes.RequireAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	})))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}
