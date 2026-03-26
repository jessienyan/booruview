package gelbooru_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func init() {
	testutil.Setup()
}

func newTestServer(handler http.HandlerFunc) (gelbooru.GelbooruClient, *httptest.Server) {
	server := httptest.NewServer(handler)
	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL
	return client, server
}

func writeJSON(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//--------------------------------------------------------

func TestSearchTags_MakesExpectedAPICall(t *testing.T) {
	called := false
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "autocomplete2", query.Get("page"))
		require.Equal(t, "foo", query.Get("term"))

		called = true
	})
	defer server.Close()

	client.SearchTags("foo")
	require.True(t, called)
}

func TestSearchTags_ParsesResponse(t *testing.T) {
	expected := []api.TagResponse{
		{Name: "foo", Type: api.Tag, Count: 1},
		{Name: "foobar", Type: api.Tag, Count: 2},
	}
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		mockResponse := []gelbooru.TagSearchResponse{
			{Type: "tag", Label: "foo", Value: "foo", Count: "1", Category: "tag"},
			{Type: "tag", Label: "foobar", Value: "foobar", Count: "2", Category: "tag"},
		}
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	actual, err := client.SearchTags("foo")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}

func TestSearchTags_ReturnsRatingSuggestions(t *testing.T) {
	expected := []api.TagResponse{
		{Name: "rating:general", Type: api.Metadata},
		{Name: "rating:questionable", Type: api.Metadata},
		{Name: "rating:sensitive", Type: api.Metadata},
		{Name: "rating:explicit", Type: api.Metadata},
	}
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	})
	defer server.Close()

	actual, err := client.SearchTags("rating:")
	require.NoError(t, err)
	require.ElementsMatch(t, expected, actual)
}

func TestSearchTags_ReturnsSortSuggestions(t *testing.T) {
	expected := []api.TagResponse{
		{Name: "sort:random", Type: api.Unknown},
		{Name: "sort:score", Type: api.Unknown},
		{Name: "sort:id", Type: api.Unknown},
		{Name: "sort:rating", Type: api.Unknown},
		{Name: "sort:user", Type: api.Unknown},
		{Name: "sort:height", Type: api.Unknown},
		{Name: "sort:width", Type: api.Unknown},
		{Name: "sort:source", Type: api.Unknown},
		{Name: "sort:updated", Type: api.Unknown},
	}
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	})
	defer server.Close()

	actual, err := client.SearchTags("sort:")
	require.NoError(t, err)
	require.ElementsMatch(t, expected, actual)
}

func TestSearchTags_OverridesFilterSearches(t *testing.T) {
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "should not be called")
	})
	defer server.Close()

	for key := range gelbooru.SearchFilters {
		_, err := client.SearchTags("rating:" + key)
		require.NoError(t, err)
	}
}

//--------------------------------------------------------

func TestListPosts_MakesExpectedAPICall(t *testing.T) {
	called := false
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "post", query.Get("s"))
		require.Equal(t, "index", query.Get("q"))
		require.Equal(t, "100", query.Get("limit"))
		require.Equal(t, "1", query.Get("json"))
		require.Equal(t, "a&amp;b -c sort:score", query.Get("tags"))
		require.Equal(t, "0", query.Get("pid"))

		called = true
	})
	defer server.Close()

	client.ListPosts("a&b -c sort:score", 1)
	require.True(t, called)
}

func TestListPosts_ParsesResponse(t *testing.T) {
	expected := gelbooru.PostList{
		TotalCount: 500,
		Posts: []api.PostResponse{
			{
				Id:                 1,
				CreatedAtTimestamp: time.Date(2026, 2, 16, 1, 23, 45, 0, time.UTC).Unix(),
				Score:              100,
				Rating:             "general",
				SourceUrl:          "https://example.com/sauce",
				Uploader:           "user1",
				UploaderUrl:        "https://gelbooru.com/index.php?page=account&s=profile&id=123",
				Tags:               []string{"tag1", "tag2", "rating:general"},
				ImageUrl:           "https://example.com/full.jpg",
				Width:              800,
				Height:             600,
				ThumbnailUrl:       "https://example.com/preview.jpg",
				ThumbnailWidth:     200,
				ThumbnailHeight:    150,
				LowResUrl:          "https://example.com/sample.jpg",
				LowResWidth:        400,
				LowResHeight:       300,
			},
		}}
	mockResponse := gelbooru.FullPostResponse{
		Attributes: gelbooru.ResultListInfo{
			Limit:  100,
			Offset: 0,
			Count:  500,
		},
		Post: []gelbooru.PostResponse{
			{
				Id:            1,
				CreatedAt:     "Mon Feb 16 01:23:45 -0000 2026",
				Score:         100,
				Rating:        "general",
				SourceUrl:     "https://example.com/sauce",
				Uploader:      "user1",
				UploaderId:    123,
				Tags:          "tag1 tag2",
				ImageUrl:      "https://example.com/full.jpg",
				Width:         800,
				Height:        600,
				PreviewUrl:    "https://example.com/preview.jpg",
				PreviewWidth:  200,
				PreviewHeight: 150,
				SampleUrl:     "https://example.com/sample.jpg",
				SampleWidth:   400,
				SampleHeight:  300,
			},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	actual, err := client.ListPosts("", 1)
	require.NoError(t, err)
	require.Equal(t, expected, *actual)
}

func TestListPosts_IncludesRatingAsMetadataTag(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Post: []gelbooru.PostResponse{
			{
				Rating: "sensitive",
				Tags:   "some_tag",
			},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	postList, _ := client.ListPosts("", 1)

	actual := postList.Posts[0].Tags
	require.Equal(t, []string{"some_tag", "rating:sensitive"}, actual)
}

func TestListPosts_UnescapesHTMLTags(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Post: []gelbooru.PostResponse{
			{
				Rating: "general",
				Tags:   "escaped&amp;tag",
			},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	postList, _ := client.ListPosts("", 1)

	tags := postList.Posts[0].Tags
	require.Equal(t, "escaped&tag", tags[0])
}

func TestListPosts_RewritesVideoCDN(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Post: []gelbooru.PostResponse{
			{
				Rating:   "general",
				ImageUrl: "https://video-cdn3.gelbooru.com/video.mp4",
			},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	postList, _ := client.ListPosts("", 1)

	require.Equal(t, 1, len(postList.Posts))
	expectedURL := "https://video-cdn4.gelbooru.com/video.mp4"
	require.Equal(t, expectedURL, postList.Posts[0].ImageUrl)
}

func TestListPosts_RewritesProxyURL(t *testing.T) {
	api.UseMediaProxy = true
	api.MediaProxyHost = "https://media.proxy"
	defer func() { api.UseMediaProxy = false }()

	mockResponse := gelbooru.FullPostResponse{
		Post: []gelbooru.PostResponse{
			{
				ImageUrl:   "https://example.com/full.jpg",
				PreviewUrl: "https://example.com/preview.jpg",
				SampleUrl:  "https://example.com/sample.jpg",
			},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	postList, _ := client.ListPosts("", 1)
	post := postList.Posts[0]

	require.Equal(
		t,
		"https://media.proxy/?to="+url.PathEscape("https://example.com/full.jpg"),
		post.ImageUrl,
	)
	require.Equal(
		t,
		"https://media.proxy/?to="+url.PathEscape("https://example.com/preview.jpg"),
		post.ThumbnailUrl,
	)
	require.Equal(
		t,
		"https://media.proxy/?to="+url.PathEscape("https://example.com/sample.jpg"),
		post.LowResUrl,
	)
}

//--------------------------------------------------------

func TestListTags_MakesExpectedAPICall(t *testing.T) {
	called := false
	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "tag", query.Get("s"))
		require.Equal(t, "test_tag artist_name", query.Get("names"))

		called = true
	})
	defer server.Close()

	client.ListTags("test_tag artist_name")
	require.True(t, called)
}

func TestListTags_ParsesResponse(t *testing.T) {
	expected := []api.TagResponse{
		{Name: "test_tag", Count: 100, Type: api.Tag},
		{Name: "artist_name", Count: 50, Type: api.Artist},
	}
	mockResponse := gelbooru.FullTagInfoResponse{
		Attributes: gelbooru.ResultListInfo{
			Limit:  100,
			Offset: 0,
			Count:  2,
		},
		Tag: []gelbooru.TagInfo{
			{Id: 1, Name: "test_tag", Count: 100, Type: 0, Ambiguous: 0},
			{Id: 2, Name: "artist_name", Count: 50, Type: 1, Ambiguous: 0},
		},
	}

	client, server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, mockResponse, 200)
	})
	defer server.Close()

	actual, err := client.ListTags("test_tag artist_name")
	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
