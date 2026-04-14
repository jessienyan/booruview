package routes

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/valkey-io/valkey-go"
)

type PostsResponse struct {
	CountPerPage int          `json:"count_per_page"`
	TotalCount   int          `json:"total_count"`
	Results      api.PostList `json:"results"`
}

type PostsHandler struct {
	Client gelbooru.GelbooruClient
}

func (h PostsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// NOTE: post rate limiting happens after checking the cache. The cost increases
	// if there's a cache miss
	if err := req.ParseForm(); err != nil {
		err = errors.Wrap(err, "failed to parse form")
		respondWithInternalError(w, err)
		return
	}

	pageVal := req.FormValue("page")
	if pageVal == "" {
		pageVal = "1"
	}

	page, err := strconv.Atoi(pageVal)
	if err != nil || page < 1 {
		respondWithBadRequest(w, "invalid page number")
		return
	}

	if page > 200 {
		respondWithBadRequest(w, "results past page 200 are blocked by gelbooru")
		return
	}

	resp := PostsResponse{
		Results: api.PostList{},
	}
	tags := api.CleanTagList(req.Form["q"])
	query := strings.Join(tags, " ")

	cached, err := GetCachedPosts(query, page)
	if err != nil {
		err = errors.Wrap(err, "failed to get cached posts")
		respondWithInternalError(w, err)
		return
	}

	// Cache hit
	if cached != nil {
		if isRateLimited(w, req, postApiCostIfCacheHit) {
			return
		}

		respondJson(w, 200, api.DecompressData(cached))
		return
	}

	// Cache miss, check with increased rate limit cost
	if isRateLimited(w, req, postApiCostIfCacheMiss) {
		return
	}

	results, err := h.Client.ListPosts(query, page)
	if err != nil {
		if errors.As(err, &gelbooru.GelbooruError{}) {
			respondWithGelbooruUnavailable(w)
			return
		}

		err = errors.Wrap(err, "failed to list posts")
		respondWithInternalError(w, err)
		return
	}

	resp.CountPerPage = gelbooru.PostsPerPage
	resp.TotalCount = results.TotalCount
	resp.Results = results.Posts

	respData := respondJson(w, http.StatusOK, resp)
	WritePostsToCache(query, page, respData)
}

func GetCachedPosts(tags string, page int) ([]byte, error) {
	vk := api.Valkey()
	cached := vk.Do(context.Background(),
		vk.B().
			Get().
			Key(gelbooru.PostCacheKey(tags, page)).
			Build(),
	)

	if err := cached.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	data, err := cached.AsBytes()
	if err != nil {
		return nil, err
	}

	return data, nil
}

func WritePostsToCache(query string, afterId int, data []byte) error {
	vk := api.Valkey()
	compressed, err := api.CompressData(data)
	if err != nil {
		return err
	}

	return vk.Do(context.Background(),
		vk.B().
			Setex().
			Key(gelbooru.PostCacheKey(query, afterId)).
			Seconds(api.PostTtl).
			Value(string(compressed)).
			Build(),
	).Error()
}
