package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.Setup()
}

func requirePostResponseEqual(t *testing.T, expected gelbooru.PostList, actual string) {
	fullExpected := PostsResponse{
		CountPerPage: gelbooru.PostsPerPage,
		TotalCount:   expected.TotalCount,
		Results:      expected.Posts,
	}

	require.JSONEq(t, string(testutil.MustMarshalJSON(fullExpected)), actual)
}

func TestPostsHandler_EmptyQuery(t *testing.T) {
	api.FlushRateLimits()
	expected := gelbooru.PostList{
		TotalCount: 0,
		Posts:      []api.PostResponse{},
	}
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "", 1).Return(&expected, nil)

	req := httptest.NewRequest("POST", "/posts", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertExpectations(t)
}

func TestPostsHandler_MultipleQueries(t *testing.T) {
	api.FlushRateLimits()
	expected := gelbooru.PostList{
		TotalCount: 2,
		Posts: []api.PostResponse{
			{Id: 1, Tags: []string{"test1", "test2"}},
			{Id: 2, Tags: []string{"test1", "test2"}},
		},
	}
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "test1 test2", 1).Return(&expected, nil)

	req := httptest.NewRequest("POST", "/posts?q=test1&q=test2", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertExpectations(t)
}

func TestPostsHandler_DefaultPage(t *testing.T) {
	api.FlushRateLimits()
	expected := gelbooru.PostList{}
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "test", 1).Return(&expected, nil)

	req := httptest.NewRequest("POST", "/posts", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertExpectations(t)
}

func TestPostsHandler_ValidQuery(t *testing.T) {
	api.FlushRateLimits()
	expected := gelbooru.PostList{
		TotalCount: 3,
		Posts: []api.PostResponse{
			{Id: 1, Tags: []string{"test"}},
			{Id: 2, Tags: []string{"test"}},
			{Id: 3, Tags: []string{"test"}},
		},
	}
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "test", 1).Return(&expected, nil)

	req := httptest.NewRequest("POST", "/posts?q=test&page=1", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertExpectations(t)
}

func TestPostsHandler_InvalidPage(t *testing.T) {
	api.FlushRateLimits()
	req := httptest.NewRequest("POST", "/posts?page=invalid", nil)
	rec := httptest.NewRecorder()

	PostsHandler{}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "invalid page")
}

func TestPostsHandler_PageExceedsLimit(t *testing.T) {
	api.FlushRateLimits()
	req := httptest.NewRequest("POST", "/posts?page=201", nil)
	rec := httptest.NewRecorder()

	PostsHandler{}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "page 200")
}

func TestPostsHandler_CacheMiss(t *testing.T) {
	api.FlushRateLimits()
	expected := gelbooru.PostList{
		TotalCount: 5,
		Posts: []api.PostResponse{
			{Id: 1, Tags: []string{"test"}},
		},
	}
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "test", 1).Return(&expected, nil)

	req := httptest.NewRequest("POST", "/posts?q=test&page=1", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertExpectations(t)
}

func TestPostsHandler_CacheHit(t *testing.T) {
	api.FlushRateLimits()
	client := &testutil.MockGelbooruClient{}
	expected := gelbooru.PostList{
		TotalCount: 1,
		Posts:      []api.PostResponse{{Id: 1, Tags: []string{"test"}}},
	}

	rec := httptest.NewRecorder()
	data := respondJson(rec, http.StatusOK, expected)
	writePostsToCache("test", 1, data)

	req := httptest.NewRequest("POST", "/posts?q=test&page=1", nil)
	rec2 := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec2, req)

	require.Equal(t, http.StatusOK, rec2.Code)
	requirePostResponseEqual(t, expected, rec.Body.String())
	client.AssertNotCalled(t, "ListPosts")
}

func TestPostsHandler_GelbooruUnavailable(t *testing.T) {
	api.FlushRateLimits()
	client := &testutil.MockGelbooruClient{}

	testutil.FlushCache()
	client.On("ListPosts", "", 1).Return(nil, gelbooru.GelbooruError{})

	req := httptest.NewRequest("POST", "/posts", nil)
	rec := httptest.NewRecorder()

	PostsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}
