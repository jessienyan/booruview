package routes

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/rs/zerolog/log"
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

	isNaughty := api.NaughtyFingerprints[req.Header.Get("Ja4h")]
	resp := PostsResponse{
		Results: []api.PostResponse{},
	}

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

	tags := req.Form["q"]

	// Clean up the query so we're left with a sorted list of unique tags
	normalized := slices.DeleteFunc(
		tags,
		func(s string) bool {
			return len(s) == 0
		},
	)
	slices.Sort(normalized)

	query := strings.Join(normalized, " ")
	cached, err := getCachedPosts(query, page)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	// Nice users get a cache :)
	if !isNaughty && cached != nil {
		if isRateLimited(w, req, postApiCostIfCacheHit) {
			return
		}

		respondJson(w, 200, api.DecompressData(cached))
		return
	}

	if isRateLimited(w, req, postApiCostIfCacheMiss) {
		return
	}

	var client gelbooru.Client

	if isNaughty {
		// Naughty users get fake data and latency :)
		client = gelbooru.NewFakeClient()
		fakeLatency := time.Duration(rand.Int()%2_000+500) * time.Millisecond
		log.Info().Float64("latency", fakeLatency.Seconds()).Str("ja4h", req.Header.Get("Ja4h")).Msg("sending client coal :)")
		time.Sleep(fakeLatency)
	} else {
		client = gelbooru.NewClient()
	}

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

	if !isNaughty {
		writePostsToCache(query, page, respData)
	}
}

func getCachedPosts(tags string, page int) ([]byte, error) {
	vc := api.Valkey()
	cached := vc.Do(context.Background(),
		vc.B().
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
	vc := api.Valkey()
	buf := bytes.Buffer{}
	if err := api.CompressData(&buf, data); err != nil {
		return err
	}

	return vc.Do(context.Background(),
		vc.B().
			Setex().
			Key(gelbooru.PostCacheKey(query, afterId)).
			Seconds(api.PostTtl).
			Value(buf.String()).
			Build(),
	).Error()
}
