package routes_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"codeberg.org/jessienyan/booruview/routes"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.Setup()
	api.InitUserDatabase()
}

func TestRegisterHandler_Success(t *testing.T) {
	testutil.Flush()

	params := routes.CreateUserParams{
		Username: "regtestuser",
		Password: "password",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.RegisterResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	parsedToken, err := api.ParseAuthToken(response.AuthToken)
	require.NoError(t, err)

	db := models.New(api.UserDB())
	user, err := db.GetUser(t.Context(), params.Username)
	require.NoError(t, err)
	require.Equal(t, parsedToken.UserID, user.ID)
}

func TestRegisterHandler_SpecialCharacters(t *testing.T) {
	testutil.Flush()

	params := routes.CreateUserParams{
		Username: "$pecial",
		Password: "password",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_UsernameTooShort(t *testing.T) {
	testutil.Flush()

	params := routes.CreateUserParams{
		Username: "a",
		Password: "password",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_UsernameTooLong(t *testing.T) {
	testutil.Flush()

	params := routes.CreateUserParams{
		Username: "thisusernameistoolong",
		Password: "password",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_PasswordTooShort(t *testing.T) {
	testutil.Flush()

	params := routes.CreateUserParams{
		Username: "blah",
		Password: "abc",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestRegisterHandler_UserExists(t *testing.T) {
	testutil.Flush()

	testutil.CreateUser("existinguser", "password4")

	params := routes.CreateUserParams{
		Username: "EXISTINGUSER",
		Password: "password",
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	routes.RegisterHandler(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "username is already taken")
}
