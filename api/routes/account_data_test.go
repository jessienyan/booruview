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
		Date: time.Now(),
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

	post := api.PostResponse{Id: 12345}
	params := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
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

	require.Len(t, response.FavoritePosts, 1)
	require.Equal(t, post.Id, response.FavoritePosts[0].Id)
}

func TestAccountDataPatchHandler_AddFavoriteTags(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	tag := api.TagResponse{Name: "test_tag", Count: 100, Type: api.Tag}
	params := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
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

	require.Len(t, response.FavoriteTags, 1)
	require.Equal(t, tag.Name, response.FavoriteTags[0].Name)
}

func TestAccountDataPatchHandler_AddBlacklist(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	tag := api.TagResponse{Name: "blacklisted", Count: 50, Type: api.Tag}
	params := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
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

	require.Len(t, response.Blacklist, 1)
	require.Equal(t, tag.Name, response.Blacklist[0].Name)
}

func TestAccountDataPatchHandler_AddSearchHistory(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	now := time.Now()
	params := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			SearchHistory: models.SearchHistoryList{
				{
					Date: now,
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

	require.Len(t, response.SearchHistory, 1)
	require.Equal(t, now.Unix(), response.SearchHistory[0].Date.Unix())
	require.Len(t, response.SearchHistory[0].Query.Include, 1)
	require.Len(t, response.SearchHistory[0].Query.Exclude, 1)
}

func TestAccountDataPatchHandler_RemoveFavoritePosts(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add a post
	addParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoritePosts: api.PostList{{Id: 111}, {Id: 222}, {Id: 333}},
		},
	}
	addBody := testutil.MustMarshalJSON(addParams)
	addReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := callHandler(routes.AccountDataPatchHandler, addReq, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec.Code)

	// Now remove one post
	removeParams := routes.AccountDataPatchParams{
		Remove: &routes.RemoveAccountData{
			FavoritePostIDs: []int{222},
		},
	}
	removeBody := testutil.MustMarshalJSON(removeParams)
	removeReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(removeBody))
	removeReq.Header.Set("Content-Type", "application/json")
	removeRec := callHandler(routes.AccountDataPatchHandler, removeReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, removeRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(removeRec.Body.Bytes(), &response)

	require.Len(t, response.FavoritePosts, 2)
	require.NotContains(t, []int{response.FavoritePosts[0].Id, response.FavoritePosts[1].Id}, 222)
}

func TestAccountDataPatchHandler_RemoveFavoriteTags(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add tags
	addParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoriteTags: api.TagList{
				{Name: "tag1", Count: 10, Type: api.Tag},
				{Name: "tag2", Count: 20, Type: api.Tag},
				{Name: "tag3", Count: 30, Type: api.Tag},
			},
		},
	}
	addBody := testutil.MustMarshalJSON(addParams)
	addReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := callHandler(routes.AccountDataPatchHandler, addReq, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec.Code)

	// Now remove one tag
	removeParams := routes.AccountDataPatchParams{
		Remove: &routes.RemoveAccountData{
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

	require.Len(t, response.FavoriteTags, 2)
	tagNames := []string{response.FavoriteTags[0].Name, response.FavoriteTags[1].Name}
	require.NotContains(t, tagNames, "tag2")
}

func TestAccountDataPatchHandler_RemoveBlacklist(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add to blacklist
	addParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			Blacklist: api.TagList{
				{Name: "bad1", Count: 10, Type: api.Tag},
				{Name: "bad2", Count: 20, Type: api.Tag},
				{Name: "bad3", Count: 30, Type: api.Tag},
			},
		},
	}
	addBody := testutil.MustMarshalJSON(addParams)
	addReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := callHandler(routes.AccountDataPatchHandler, addReq, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec.Code)

	// Now remove from blacklist
	removeParams := routes.AccountDataPatchParams{
		Remove: &routes.RemoveAccountData{
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

	require.Len(t, response.Blacklist, 2)
	blacklistNames := []string{response.Blacklist[0].Name, response.Blacklist[1].Name}
	require.NotContains(t, blacklistNames, "bad2")
}

func TestAccountDataPatchHandler_RemoveSearchHistory(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add search history
	now := time.Now()
	addParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			SearchHistory: models.SearchHistoryList{
				{
					Date: now,
					Query: models.SearchQuery{
						Include: api.TagList{api.TagResponse{Name: "query1", Count: 10, Type: api.Tag}},
						Exclude: api.TagList{},
					},
				},
				{
					Date: now,
					Query: models.SearchQuery{
						Include: api.TagList{api.TagResponse{Name: "query2", Count: 20, Type: api.Tag}},
						Exclude: api.TagList{},
					},
				},
			},
		},
	}
	addBody := testutil.MustMarshalJSON(addParams)
	addReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := callHandler(routes.AccountDataPatchHandler, addReq, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec.Code)

	// Now remove one search history entry
	removeParams := routes.AccountDataPatchParams{
		Remove: &routes.RemoveAccountData{
			SearchQueries: []models.SearchQuery{
				{
					Include: api.TagList{api.TagResponse{Name: "query2", Count: 20, Type: api.Tag}},
					Exclude: api.TagList{},
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

	require.Len(t, response.SearchHistory, 1)
	require.Equal(t, "query1", response.SearchHistory[0].Query.Include[0].Name)
}

func TestAccountDataPatchHandler_MixedAddAndRemove(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add some data
	addParams1 := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoriteTags: api.TagList{
				{Name: "tag1", Count: 10, Type: api.Tag},
				{Name: "tag2", Count: 20, Type: api.Tag},
			},
		},
	}
	addBody1 := testutil.MustMarshalJSON(addParams1)
	addReq1 := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody1))
	addReq1.Header.Set("Content-Type", "application/json")
	addRec1 := callHandler(routes.AccountDataPatchHandler, addReq1, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec1.Code)

	// Mix add and remove in one request
	mixedParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoriteTags: api.TagList{{Name: "tag3", Count: 30, Type: api.Tag}},
		},
		Remove: &routes.RemoveAccountData{
			FavoriteTagNames: []string{"tag1"},
		},
	}
	mixedBody := testutil.MustMarshalJSON(mixedParams)
	mixedReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(mixedBody))
	mixedReq.Header.Set("Content-Type", "application/json")
	mixedRec := callHandler(routes.AccountDataPatchHandler, mixedReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, mixedRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(mixedRec.Body.Bytes(), &response)

	require.Len(t, response.FavoriteTags, 2)
	tagNames := []string{response.FavoriteTags[0].Name, response.FavoriteTags[1].Name}
	require.Contains(t, tagNames, "tag2")
	require.Contains(t, tagNames, "tag3")
	require.NotContains(t, tagNames, "tag1")
}

func TestAccountDataPatchHandler_InvalidContentType(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	params := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{},
	}
	body := testutil.MustMarshalJSON(params)

	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	// Missing Content-Type header
	rec := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestAccountDataPutHandler_ReplaceFavoritePosts(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// First add some posts
	addParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoritePosts: api.PostList{{Id: 111}, {Id: 222}},
		},
	}
	addBody := testutil.MustMarshalJSON(addParams)
	addReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody))
	addReq.Header.Set("Content-Type", "application/json")
	addRec := callHandler(routes.AccountDataPatchHandler, addReq, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec.Code)

	// Now replace all posts with PUT
	putParams := routes.AccountDataPutParams{
		UserDataJSON: models.UserDataJSON{
			FavoritePosts: api.PostList{{Id: 333}, {Id: 444}},
		},
	}
	putBody := testutil.MustMarshalJSON(putParams)
	putReq := httptest.NewRequest("PUT", "/api/account/data", bytes.NewReader(putBody))
	putReq.Header.Set("Content-Type", "application/json")
	putRec := callHandler(routes.AccountDataPutHandler, putReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, putRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(putRec.Body.Bytes(), &response)

	require.Len(t, response.FavoritePosts, 2)
	require.Equal(t, 333, response.FavoritePosts[0].Id)
	require.Equal(t, 444, response.FavoritePosts[1].Id)
}

func TestAccountDataPatchHandler_EmptyAdd(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	// Add some data first
	addParams1 := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{
			FavoriteTags: api.TagList{{Name: "tag1", Count: 10, Type: api.Tag}},
		},
	}
	addBody1 := testutil.MustMarshalJSON(addParams1)
	addReq1 := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(addBody1))
	addReq1.Header.Set("Content-Type", "application/json")
	addRec1 := callHandler(routes.AccountDataPatchHandler, addReq1, accountDataAuthToken)
	require.Equal(t, http.StatusOK, addRec1.Code)

	// Now send empty add - should return current state without changes
	emptyParams := routes.AccountDataPatchParams{
		Add: &routes.AddAccountData{},
	}
	emptyBody := testutil.MustMarshalJSON(emptyParams)
	emptyReq := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(emptyBody))
	emptyReq.Header.Set("Content-Type", "application/json")
	emptyRec := callHandler(routes.AccountDataPatchHandler, emptyReq, accountDataAuthToken)

	require.Equal(t, http.StatusOK, emptyRec.Code)

	var response routes.AccountDataResponse
	testutil.MustUnmarshalJSON(emptyRec.Body.Bytes(), &response)

	// Data should be unchanged
	require.Len(t, response.FavoriteTags, 1)
	require.Equal(t, "tag1", response.FavoriteTags[0].Name)
}
