package routes

import (
	"context"
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

func TestTagSearchHandler_EmptyQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/tagsearch?q=", nil)
	rec := httptest.NewRecorder()

	TagSearchHandler{}.ServeHTTP(rec, req)

	require.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestTagSearchHandler_CacheHit(t *testing.T) {
	client := &testutil.MockGelbooruClient{}
	req := httptest.NewRequest("GET", "/tagsearch?q=sky", nil)
	rec := httptest.NewRecorder()

	writeTagSearchToCache("sky", []byte("foo"))
	TagSearchHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, rec.Body.Bytes(), []byte("foo"))
	require.Equal(t, rec.Code, http.StatusOK)
	require.Equal(t, rec.Header().Get("Content-Type"), "application/json")
	client.AssertNotCalled(t, "SearchTags")
}

func TestTagSearchHandler_NoResults(t *testing.T) {
	expected := []api.TagResponse{}
	client := &testutil.MockGelbooruClient{}
	req := httptest.NewRequest("GET", "/tagsearch?q=test", nil)
	rec := httptest.NewRecorder()

	testutil.FlushCache()
	client.On("SearchTags", "test").Return(expected, nil)
	TagSearchHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, rec.Code, http.StatusOK)
	require.Equal(t, rec.Header().Get("Content-Type"), "application/json")
	require.JSONEq(t, rec.Body.String(), `{"results": []}`)
	client.AssertExpectations(t)
}

func TestTagSearchHandler_SomeResults(t *testing.T) {
	results := []api.TagResponse{{Name: "test", Count: 0, Type: api.Tag}}
	client := &testutil.MockGelbooruClient{}
	req := httptest.NewRequest("GET", "/tagsearch?q=test", nil)
	rec := httptest.NewRecorder()

	testutil.FlushCache()
	client.On("SearchTags", "test").Return(results, nil)
	TagSearchHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, rec.Code, http.StatusOK)
	require.Equal(t, rec.Header().Get("Content-Type"), "application/json")
	require.JSONEq(t, rec.Body.String(), `{"results": [{"name": "test", "count": 0, "type": "tag"}]}`)
	client.AssertExpectations(t)

	// Check it was cached
	vk := api.Valkey()
	key := gelbooru.TagSearchCacheKey("test")
	exists, _ := vk.Do(context.Background(), vk.B().Exists().Key(key).Build()).AsBool()
	require.True(t, exists)
}

func TestTagSearchHandler_GelbooruDown(t *testing.T) {
	client := &testutil.MockGelbooruClient{}
	req := httptest.NewRequest("GET", "/tagsearch?q=test", nil)
	rec := httptest.NewRecorder()

	testutil.FlushCache()
	client.On("SearchTags", "test").Return([]api.TagResponse{}, gelbooru.GelbooruError{})
	TagSearchHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, rec.Code, http.StatusServiceUnavailable)
}
