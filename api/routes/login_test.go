package routes_test

import (
	"bytes"
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
	loginTestUser     models.Users
	loginTestUsername = "logintest"
	loginTestPassword = "pass123"
)

func init() {
	testutil.Setup()
	api.InitUserDatabase()
	loginTestUser, _ = testutil.CreateUser(loginTestUsername, loginTestPassword)
}

func TestLoginHandler_UserDoesntExist(t *testing.T) {
	testutil.Flush()

	params := routes.LoginParams{
		Username: "lkjsdfgoijwefa",
		Password: "wrong",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.LoginHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_IncorrectPassword(t *testing.T) {
	testutil.Flush()

	params := routes.LoginParams{
		Username: loginTestUsername,
		Password: "wrongagain",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.LoginHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	testutil.Flush()

	params := routes.LoginParams{
		Username: loginTestUsername,
		Password: loginTestPassword,
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.LoginHandler(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.LoginResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.NotEmpty(t, response.AuthToken)
	require.Equal(t, loginTestUser.Username, response.Username)
}

func TestLoginHandler_UpdatesLastLogin(t *testing.T) {
	testutil.Flush()

	params := routes.LoginParams{
		Username: loginTestUsername,
		Password: loginTestPassword,
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	earlier := time.Now()
	routes.LoginHandler(rec, req)

	db := models.New(api.UserDB())
	updatedUser, _ := db.GetUserByID(t.Context(), loginTestUser.ID)

	require.True(t, updatedUser.LastLogin.Valid)
	require.True(t, updatedUser.LastLogin.Time.After(earlier))
}

func TestLoginHandler_RequiredFields(t *testing.T) {
	testutil.Flush()

	req := httptest.NewRequest("POST", "/login", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.LoginHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "username and password are required")
}
