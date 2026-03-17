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
	authMiddlewareUser = testutil.CreateUser(authMiddlewareUsername, authMiddlewarePassword)
}

func TestAuthMiddleware_MissingAuthHeader(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("GET", "/faked", nil)
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("GET", "/faked", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	testutil.Flush()

	token, _ := api.NewAuthToken(int(authMiddlewareUser.ID), 0)

	req := httptest.NewRequest("GET", "/faked", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_NonExistentUser(t *testing.T) {
	testutil.Flush()

	token, _ := api.NewAuthToken(123456789, 1*time.Minute)

	req := httptest.NewRequest("GET", "/faked", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusUnauthorized, rec.Code)
}

func TestAuthMiddleware_ValidTokenCallsNextHandler(t *testing.T) {
	testutil.Flush()

	token, _ := api.NewAuthToken(int(authMiddlewareUser.ID), 1*time.Minute)

	req := httptest.NewRequest("GET", "/faked", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()

	called := false
	handler := routes.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true // sanity check
		w.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.True(t, called)
}
