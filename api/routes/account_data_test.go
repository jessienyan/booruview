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
	accountDataAuthToken    string
	accountDataTestUsername = "accountdatatest"
	accountDataTestPassword = "pass123"
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
	resp := httptest.NewRecorder()
	wrappedHandler := routes.AuthMiddleware(http.HandlerFunc(handlerFunc))
	wrappedHandler.ServeHTTP(resp, req)
	return resp
}

func patchAccountData(data any) *httptest.ResponseRecorder {
	body := testutil.MustMarshalJSON(data)
	req := httptest.NewRequest("PATCH", "/api/account/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp := callHandler(routes.AccountDataPatchHandler, req, accountDataAuthToken)
	return resp
}

func TestAccountDataGetHandler_Empty(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	req := httptest.NewRequest("GET", "/api/account/data", nil)
	resp := callHandler(routes.AccountDataGetHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataResponse{
		models.UserDataJSON{
			FavoritePosts: api.PostList{},
			FavoriteTags:  api.TagList{},
			Blacklist:     api.TagList{},
			SearchHistory: models.SearchHistoryList{},
			SavedSearches: models.SearchHistoryList{},
		},
	}

	require.Equal(t, expected, actual)
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
	data.SavedSearches = models.SearchHistoryList{{
		Date: testutil.Time(),
		Query: models.SearchQuery{
			Include: api.TagList{{Name: "test5", Type: api.Tag}},
			Exclude: api.TagList{{Name: "test6", Type: api.Metadata}},
		},
	}}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	// Explicitly using JSON here so it's tested at least once
	expected := `{
	  "favorite_posts": [
		{
		  "id": 1,
		  "created_at": 0,
		  "score": 0,
		  "rating": "",
		  "source_url": "",
		  "uploader": "",
		  "uploader_url": "",
		  "tags": null,
		  "thumbnail_url": "",
		  "thumbnail_width": 0,
		  "thumbnail_height": 0,
		  "lowres_url": "",
		  "lowres_width": 0,
		  "lowres_height": 0,
		  "image_url": "example.com/blah.jpg",
		  "width": 0,
		  "height": 0
		}
	  ],
	  "favorite_tags": [
		{
		  "name": "test2",
		  "type": "character",
		  "count": 2
		}
	  ],
	  "blacklist": [
		{
		  "name": "test",
		  "type": "artist",
		  "count": 1
		}
	  ],
	  "search_history": [
		{
		  "date": "2026-04-01T01:23:45Z",
		  "query": {
			"include": [
			  {
				"name": "test3",
				"type": "tag",
				"count": 0
			  }
			],
			"exclude": [
			  {
				"name": "test4",
				"type": "metadata",
				"count": 0
			  }
			]
		  }
		}
	  ],
	  "saved_searches": [
		{
		  "date": "2026-04-01T01:23:45Z",
		  "query": {
			"include": [
			  {
				"name": "test5",
				"type": "tag",
				"count": 0
			  }
			],
			"exclude": [
			  {
				"name": "test6",
				"type": "metadata",
				"count": 0
			  }
			]
		  }
		}
	  ]
	}`

	req := httptest.NewRequest("GET", "/api/account/data", nil)
	resp := callHandler(routes.AccountDataGetHandler, req, accountDataAuthToken)

	require.Equal(t, http.StatusOK, resp.Code)
	require.JSONEq(t, expected, resp.Body.String())
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

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		FavoritePosts: api.PostList{post},
	}
	require.Equal(t, expected, actual)
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
	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)

	require.Equal(t, api.TagList{tag}, actual.FavoriteTags)
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

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		Blacklist: params.Add.Blacklist,
	}
	require.Equal(t, expected, actual)
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

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		SearchHistory: params.Add.SearchHistory,
	}
	require.Equal(t, expected, actual)
}

func TestAccountDataPatchHandler_AddSavedSearch(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	params := routes.AccountDataPatchParams{
		Add: routes.AddAccountData{
			SavedSearches: models.SearchHistoryList{
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

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		SavedSearches: params.Add.SavedSearches,
	}
	require.Equal(t, expected, actual)
}

func TestAccountDataPatchHandler_RemoveFavoritePosts(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.FavoritePosts = api.PostList{{Id: 1}, {Id: 2}, {Id: 3}}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			FavoritePostIDs: []int{2},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		FavoritePosts: api.PostList{{Id: 1}, {Id: 3}},
	}
	require.Equal(t, expected, actual)
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
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			FavoriteTagNames: []string{"tag2"},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		FavoriteTags: api.TagList{
			{Name: "tag1", Count: 1, Type: api.Tag},
			{Name: "tag3", Count: 3, Type: api.Tag},
		},
	}
	require.Equal(t, expected, actual)
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
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			BlacklistNames: []string{"bad2"},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		Blacklist: api.TagList{
			{Name: "bad1", Count: 1, Type: api.Tag},
			{Name: "bad3", Count: 3, Type: api.Tag},
		},
	}
	require.Equal(t, expected, actual)
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
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	// Now remove one search history entry
	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			SearchQueries: []models.SearchQueryNames{
				{
					Include: []string{"test1"},
					Exclude: []string{"test2"},
				},
			},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		SearchHistory: models.SearchHistoryList{data.SearchHistory[1]},
	}
	require.Equal(t, expected, actual)
}

func TestAccountDataPatchHandler_RemoveSavedSearch(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.SavedSearches = models.SearchHistoryList{
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
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	// Now remove one search history entry
	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			SavedQueries: []models.SearchQueryNames{
				{
					Include: []string{"test1"},
					Exclude: []string{"test2"},
				},
			},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{
		SavedSearches: models.SearchHistoryList{data.SavedSearches[1]},
	}
	require.Equal(t, expected, actual)
}

func TestAccountDataPatchHandler_ResponseIncludesFieldsIfEmpty(t *testing.T) {
	testutil.Flush()
	testutil.ResetUserData(accountDataTestUser.ID)

	var data models.UserDataJSON
	data.Blacklist = api.TagList{{Name: "test", Type: api.Artist, Count: 1}}
	accountDataTestUserData.Set(data)
	testutil.UpdateUserData(accountDataTestUser.ID, accountDataTestUserData)

	params := routes.AccountDataPatchParams{
		Remove: routes.RemoveAccountData{
			BlacklistNames: []string{"test"},
		},
	}

	resp := patchAccountData(params)
	require.Equal(t, http.StatusOK, resp.Code)

	var actual routes.AccountDataPatchResponse
	testutil.MustUnmarshalJSON(resp.Body.Bytes(), &actual)
	expected := routes.AccountDataPatchResponse{Blacklist: api.TagList{}}
	require.Equal(t, expected, actual)
}
