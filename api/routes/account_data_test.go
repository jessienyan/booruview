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

var (
	accountDataTestUser     models.Users
	accountDataTestUserData models.UserData
	accountDataTestUsername = "accountdatatest"
	accountDataTestPassword = "pass123"
	accountDataAuthToken    string
)

func init() {
	testutil.Setup()
	api.InitUserDatabase()
	accountDataTestUser, accountDataTestUserData = testutil.CreateUser(accountDataTestUsername, accountDataTestPassword)
	var err error
	accountDataAuthToken, err = api.NewAuthToken(int(accountDataTestUser.ID), api.AuthTokenTTL)
	if err != nil {
		panic(err)
	}
}

func callHandler(handlerFunc func(http.ResponseWriter, *http.Request), req *http.Request, token string) *httptest.ResponseRecorder {
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	wrappedHandler := routes.AuthMiddleware(http.HandlerFunc(handlerFunc))
	wrappedHandler.ServeHTTP(rec, req)
	return rec
}

func TestAccountDataGetHandler_Empty(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	req := httptest.NewRequest("GET", "/api/account/data", nil)
	rec := callHandler(routes.AccountDataGetHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.Empty(t, response.FavoritePosts)
	require.Empty(t, response.FavoriteTags)
	require.Empty(t, response.Blacklist)
	require.Empty(t, response.SearchHistory)
}

func TestAccountDataGetHandler_Success(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.Blacklist = api.TagList{{Name: "test", Type: api.Artist, Count: 1}}
	data.FavoriteTags = api.TagList{{Name: "test2", Type: api.Character, Count: 2}}
	data.FavoritePosts = api.PostList{{Id: 1, ImageUrl: "example.com/blah.jpg"}}
	data.SearchHistory = models.SearchHistoryList{{
		Date: testutil.Time(),
		Query: models.SearchQuery{
			Include: api.TagList{{Name: "test3", Type: api.Tag}},
			Exclude: api.TagList{{Name: "test4", Type: api.Metadata}},
		},
	}}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	expected := testutil.MustMarshalJSON(data)

	req := httptest.NewRequest("GET", "/api/account/data", nil)
	rec := callHandler(routes.AccountDataGetHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, string(expected), rec.Body.String())
}

func TestAccountDataPatchHandler_AddFavoritePosts(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	post := api.PostResponse{Id: 12345, Tags: []string{}}
	params := routes.AccountDataPatchParams{
		Add: routes.AddAccountData{
			FavoritePosts: api.PostList{post},
		},
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.Equal(t, api.PostList{post}, response.FavoritePosts)
}

func TestAccountDataPatchHandler_AddFavoriteTags(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	tag := api.TagResponse{Name: "test_tag", Count: 100, Type: api.Tag}
	params := routes.AccountDataPatchParams{
		Add: routes.AddAccountData{
			FavoriteTags: api.TagList{tag},
		},
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.Equal(t, api.TagList{tag}, response.FavoriteTags)
}

func TestAccountDataPatchHandler_AddBlacklist(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	tag := api.TagResponse{Name: "blacklisted", Count: 50, Type: api.Tag}
	params := routes.AccountDataPatchParams{
		Add: routes.AddAccountData{
			Blacklist: api.TagList{tag},
		},
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.Equal(t, api.TagList{tag}, response.Blacklist)
}

func TestAccountDataPatchHandler_AddSearchHistory(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	params := routes.AccountDataPatchParams{
		Add: routes.AddAccountData{
			SearchHistory: models.SearchHistoryList{
				{
					Date: testutil.Time(),
					Query: models.SearchQuery{
						Include: api.TagList{api.TagResponse{Name: "include_tag", Count: 10, Type: api.Tag}},
						Exclude: api.TagList{api.TagResponse{Name: "exclude_tag", Count: 5, Type: api.Tag}},
					},
				},
			},
		},
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, rec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(rec.Body.Bytes(), &response)

	require.Equal(t, params.Add.SearchHistory, response.SearchHistory)
}

func TestAccountDataPatchHandler_RemoveFavoritePosts(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.FavoritePosts = api.PostList{{Id: 1}, {Id: 2}, {Id: 3}}
	expected := api.PostList{{Id: 1}, {Id: 3}}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	removeParams := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			FavoritePostIDs: []int{2},
		},
	}
	removeBody := testutil.MustMarshalJSON(removeParams)
	removeReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(removeBody))
	removeReq.Header.Set("Content-Type", "application/json")
	removeRec := callHandler(routes.AccountDataPatchHandler, removeReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, removeRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(removeRec.Body.Bytes(), &response)
	require.Equal(t, expected, response.FavoritePosts)
}

func TestAccountDataPatchHandler_RemoveFavoriteTags(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.FavoriteTags = api.TagList{
		{Name: "tag1", Count: 1, Type: api.Tag},
		{Name: "tag2", Count: 2, Type: api.Tag},
		{Name: "tag3", Count: 3, Type: api.Tag},
	}
	expected := api.TagList{
		{Name: "tag1", Count: 1, Type: api.Tag},
		{Name: "tag3", Count: 3, Type: api.Tag},
	}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	removeParams := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			FavoriteTagNames: []string{"tag2"},
		},
	}
	removeBody := testutil.MustMarshalJSON(removeParams)
	removeReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(removeBody))
	removeReq.Header.Set("Content-Type", "application/json")
	removeRec := callHandler(routes.AccountDataPatchHandler, removeReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, removeRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(removeRec.Body.Bytes(), &response)

	require.Equal(t, expected, response.FavoriteTags)
}

func TestAccountDataPatchHandler_RemoveBlacklist(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.Blacklist = api.TagList{
		{Name: "bad1", Count: 1, Type: api.Tag},
		{Name: "bad2", Count: 2, Type: api.Tag},
		{Name: "bad3", Count: 3, Type: api.Tag},
	}
	expected := api.TagList{
		{Name: "bad1", Count: 1, Type: api.Tag},
		{Name: "bad3", Count: 3, Type: api.Tag},
	}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	removeParams := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			BlacklistNames: []string{"bad2"},
		},
	}
	removeBody := testutil.MustMarshalJSON(removeParams)
	removeReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(removeBody))
	removeReq.Header.Set("Content-Type", "application/json")
	removeRec := callHandler(routes.AccountDataPatchHandler, removeReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, removeRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(removeRec.Body.Bytes(), &response)

	require.Equal(t, expected, response.Blacklist)
}

func TestAccountDataPatchHandler_RemoveSearchHistory(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.SearchHistory = models.SearchHistoryList{
		{
			Date: testutil.Time(),
			Query: models.SearchQuery{
				Include: api.TagList{{Name: "test1", Type: api.Tag}},
				Exclude: api.TagList{{Name: "test2", Type: api.Metadata}},
			},
		},
		{
			Date: testutil.Time(),
			Query: models.SearchQuery{
				Include: api.TagList{{Name: "test3", Type: api.Tag}},
				Exclude: api.TagList{{Name: "test4", Type: api.Metadata}},
			},
		},
	}
	expected := models.SearchHistoryList{data.SearchHistory[1]}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	// Now remove one search history entry
	removeParams := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			SearchQueries: []models.SearchQueryNames{
				{
					Include: []string{"test1"},
					Exclude: []string{"test2"},
				},
			},
		},
	}
	removeBody := testutil.MustMarshalJSON(removeParams)
	removeReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(removeBody))
	removeReq.Header.Set("Content-Type", "application/json")
	removeRec := callHandler(routes.AccountDataPatchHandler, removeReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, removeRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(removeRec.Body.Bytes(), &response)

	require.Equal(t, expected, response.SearchHistory)
}
