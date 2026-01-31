package routes

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/valkey-io/valkey-go"
)

type PostsResponse struct {
	CountPerPage int                `json:"count_per_page"`
	TotalCount   int                `json:"total_count"`
	Results      []api.PostResponse `json:"results"`
}

func PostsHandler(w http.ResponseWriter, req *http.Request) {
	// NOTE: post rate limiting happens after checking the cache. The cost increases
	// if there's a cache miss
	if err := req.ParseForm(); err != nil {
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
		Results: []api.PostResponse{},
	}
	tags := cleanTagList(req.Form["q"])
	query := strings.Join(tags, " ")

	cached, err := getCachedPosts(query, page)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	// Update post history
	if len(tags) > 0 {
		user := getUser(req)
		if user != nil {
			data, err := user.Data.ParseJSON()
			if err != nil {
				respondWithInternalError(w, err)
				return
			}

			data.SearchHistory.Add(models.SearchHistoryEntry{
				CreatedAt: time.Now(),
				Tags:      tags,
			})

			if err := user.Data.Set(data); err != nil {
				respondWithInternalError(w, err)
				return
			}

			db := models.New(api.UserDB())
			err = db.UpdateUserData(req.Context(),
				models.UpdateUserDataParams{
					Data:   user.Data.Data,
					UserID: user.User.ID,
				})
			if err != nil {
				respondWithInternalError(w, err)
				return
			}
		}
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

	client := gelbooru.NewClient()

	results, err := client.ListPosts(query, page)
	if err != nil {
		if errors.As(err, &gelbooru.GelbooruError{}) {
			respondWithGelbooruUnavailable(w)
			return
		}

		respondWithInternalError(w, err)
		return
	}

	resp.CountPerPage = gelbooru.PostsPerPage
	resp.TotalCount = results.TotalCount
	resp.Results = results.Posts

	respData := respondJson(w, http.StatusOK, resp)
	writePostsToCache(query, page, respData)
}

func getCachedPosts(tags string, page int) ([]byte, error) {
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

func writePostsToCache(query string, afterId int, data []byte) error {
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
