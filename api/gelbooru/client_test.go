package gelbooru_test

import (
	"encoding/json"
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

func TestSearchTags_MakesExpectedAPICall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "autocomplete2", query.Get("page"))
		require.Equal(t, "foo", query.Get("term"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]gelbooru.TagSearchResponse{})
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	_, err := client.SearchTags("foo")
	require.NoError(t, err)
}

func TestSearchTags_ParsesResponse(t *testing.T) {
	mockResponse := []gelbooru.TagSearchResponse{
		{Type: "tag", Label: "foo", Value: "foo", Count: "1", Category: "tag"},
		{Type: "tag", Label: "foobar", Value: "foobar", Count: "2", Category: "tag"},
	}
	expected := []api.TagResponse{
		{Name: "foo", Type: api.Tag, Count: 1},
		{Name: "foobar", Type: api.Tag, Count: 2},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "api shouldn't be called")
	}))
	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Fail(t, "api shouldn't be called")
	}))
	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	actual, err := client.SearchTags("sort:")
	require.NoError(t, err)
	require.ElementsMatch(t, expected, actual)
}

func TestListPosts_ReturnsPaginatedPosts(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  500,
		},
		Post: []gelbooru.PostResponse{
			{
				Id:            1,
				CreatedAt:     "2024-01-01",
				Score:         100,
				Rating:        "general",
				SourceUrl:     "https://example.com/image1",
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "post", query.Get("s"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	postList, err := client.ListPosts("", 1)
	require.NoError(t, err)

	require.Equal(t, 500, postList.TotalCount)
	require.Equal(t, 1, len(postList.Posts))
	require.Equal(t, 1, postList.Posts[0].Id)
}

func TestListPosts_IncludesRatingAsMetadataTag(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  1,
		},
		Post: []gelbooru.PostResponse{
			{
				Id:       1,
				Score:    100,
				Rating:   "sensitive",
				Tags:     "tag1 tag2",
				ImageUrl: "https://example.com/full.jpg",
				Width:    800,
				Height:   600,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "post", query.Get("s"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	postList, err := client.ListPosts("", 1)
	require.NoError(t, err)

	require.Equal(t, 1, len(postList.Posts))

	tags := postList.Posts[0].Tags
	require.Equal(t, 3, len(tags))

	lastTag := tags[len(tags)-1]
	require.Equal(t, "rating:sensitive", lastTag)
}

func TestListPosts_UnescapesHTMLTags(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  1,
		},
		Post: []gelbooru.PostResponse{
			{
				Id:       1,
				Score:    100,
				Rating:   "general",
				Tags:     "tag_1&amp;tag_2",
				ImageUrl: "https://example.com/full.jpg",
				Width:    800,
				Height:   600,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "post", query.Get("s"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	postList, err := client.ListPosts("", 1)
	require.NoError(t, err)

	require.Equal(t, 1, len(postList.Posts))

	tags := postList.Posts[0].Tags
	require.Equal(t, 2, len(tags))

	expectedTag := "tag_1&tag_2"
	require.Equal(t, expectedTag, tags[0])
}

func TestListPosts_RewritesVideoCDN3(t *testing.T) {
	mockResponse := gelbooru.FullPostResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  1,
		},
		Post: []gelbooru.PostResponse{
			{
				Id:       1,
				Score:    100,
				Rating:   "general",
				Tags:     "test",
				ImageUrl: "https://video-cdn3.gelbooru.com/video.mp4",
				Width:    800,
				Height:   600,
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "post", query.Get("s"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	postList, err := client.ListPosts("", 1)
	require.NoError(t, err)

	require.Equal(t, 1, len(postList.Posts))

	expectedURL := "https://video-cdn4.gelbooru.com/video.mp4"
	require.Equal(t, expectedURL, postList.Posts[0].ImageUrl)
}

func TestListTags_ReturnsCorrectTagInfo(t *testing.T) {
	mockResponse := gelbooru.FullTagInfoResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  2,
		},
		Tag: []gelbooru.TagInfo{
			{Id: 1, Name: "test_tag", Count: 100, Type: 0, Ambiguous: 0},
			{Id: 2, Name: "artist_name", Count: 50, Type: 1, Ambiguous: 0},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "tag", query.Get("s"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	tags, err := client.ListTags("test_tag artist_name")
	require.NoError(t, err)

	require.Equal(t, 2, len(tags))

	require.Equal(t, "test_tag", tags[0].Name)
	require.Equal(t, 100, tags[0].Count)
	require.Equal(t, api.Tag, tags[0].Type)

	require.Equal(t, "artist_name", tags[1].Name)
	require.Equal(t, api.Artist, tags[1].Type)
}

func TestListTags_SkipsEmptyTags(t *testing.T) {
	mockResponse := gelbooru.FullTagInfoResponse{
		Attributes: struct {
			Limit  int
			Offset int
			Count  int
		}{
			Limit:  100,
			Offset: 0,
			Count:  2,
		},
		Tag: []gelbooru.TagInfo{
			{Id: 1, Name: "", Count: 0, Type: 0, Ambiguous: 0},
			{Id: 2, Name: "valid_tag", Count: 50, Type: 0, Ambiguous: 0},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		require.Equal(t, "dapi", query.Get("page"))
		require.Equal(t, "tag", query.Get("s"))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := gelbooru.NewClient(server.Client())
	client.ApiUrl = server.URL

	tags, err := client.ListTags("empty_tag valid_tag")
	require.NoError(t, err)

	require.Equal(t, 2, len(tags))
	require.Equal(t, "", tags[0].Name)
	require.Equal(t, "valid_tag", tags[1].Name)
}
