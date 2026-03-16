package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/routes"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.Setup()
}

func TestTagsHandler_EmptyQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/tags?t=", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"results": []}`, rec.Body.String())
}

func TestTagsHandler_TooManyTags(t *testing.T) {
	query := ""
	for i := range routes.TagLimit + 1 {
		query += "t=" + fmt.Sprintf("tag%d", i) + "&"
	}
	// Strip trailing &
	query = query[:len(query)-1]

	req := httptest.NewRequest("GET", "/tags?"+query, nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Contains(t, rec.Body.String(), "limit of 100 tags")
}

func TestTagsHandler_CacheHit(t *testing.T) {
	client := &testutil.MockGelbooruClient{}
	tag := api.TagResponse{Name: "test", Type: api.Tag}

	routes.WriteCachedTags([]api.TagResponse{tag})

	req := httptest.NewRequest("GET", "/tags?t=test", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"results": [{"name": "test", "count": 0, "type": "tag"}]}`, rec.Body.String())
	client.AssertNotCalled(t, "ListTags")
}

func TestTagsHandler_CacheMiss(t *testing.T) {
	expected := []api.TagResponse{{Name: "test", Type: api.Tag}}
	client := &testutil.MockGelbooruClient{}

	testutil.Flush()
	client.On("ListTags", "test").Return(expected, nil)

	req := httptest.NewRequest("GET", "/tags?t=test", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"results": [{"name": "test", "count": 0, "type": "tag"}]}`, rec.Body.String())
	client.AssertExpectations(t)
}

func TestTagsHandler_PartialCacheHit(t *testing.T) {
	client := &testutil.MockGelbooruClient{}
	cachedTag := api.TagResponse{Name: "cached", Count: 1, Type: api.Tag}
	newTag := api.TagResponse{Name: "new", Count: 2, Type: api.Artist}

	testutil.Flush()
	routes.WriteCachedTags([]api.TagResponse{cachedTag})
	client.On("ListTags", "new").Return([]api.TagResponse{newTag}, nil)

	req := httptest.NewRequest("GET", "/tags?t=cached&t=new", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"results": [{"name": "cached", "count": 1, "type": "tag"}, {"name": "new", "count": 2, "type": "artist"}]}`, rec.Body.String())
	client.AssertExpectations(t)
}

func TestTagsHandler_GelbooruUnavailable(t *testing.T) {
	client := &testutil.MockGelbooruClient{}

	testutil.Flush()
	client.On("ListTags", "test").Return([]api.TagResponse{}, gelbooru.GelbooruError{})

	req := httptest.NewRequest("GET", "/tags?t=test", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestTagsHandler_StripsHyphen(t *testing.T) {
	expected := []api.TagResponse{{Name: "stripme", Type: api.Tag}}
	client := &testutil.MockGelbooruClient{}

	testutil.Flush()
	client.On("ListTags", "stripme").Return(expected, nil)

	req := httptest.NewRequest("GET", "/tags?t=-stripme", nil)
	rec := httptest.NewRecorder()

	routes.TagsHandler{Client: client}.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"results": [{"name": "stripme", "count": 0, "type": "tag"}]}`, rec.Body.String())
	client.AssertExpectations(t)
}
