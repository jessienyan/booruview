package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	api "github.com/jessienyan/booruview"
	"github.com/jessienyan/booruview/gelbooru"
	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
)

const (
	postApiCostIfCacheHit = 1

	// 5 seems a bit aggressive but the token refill is also very generous.
	// Swiping page after page on my phone I was able to do 20 pages before
	// hitting the limit (with Gelbooru latency around 700ms)
	postApiCostIfCacheMiss = 5
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

	pageVal := req.FormValue("page")
	if pageVal == "" {
		pageVal = "1"
	}

	page, err := strconv.Atoi(pageVal)
	if err != nil || page < 1 {
		handle400Error(w, "invalid page number")
		return
	}

	// Clean up the query so we're left with a sorted list of unique tags
	normalized := slices.DeleteFunc(
		strings.Split(req.FormValue("q"), ","),
		func(s string) bool {
			return len(s) == 0
		},
	)
	slices.Sort(normalized)

	query := strings.Join(normalized, " ")
	cached, err := getCachedPosts(query, page)
	if err != nil {
		handleError(w, err)
		return
	}

	// Nice users get a cache :)
	if !isNaughty && cached != nil {
		if isRateLimited(w, req, postApiCostIfCacheHit) {
			return
		}

		api.DecompressData(w, cached)
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
		if _, ok := err.(gelbooru.GelbooruError); ok {
			handle400Error(w, "Gelbooru is not available right now")
			return
		}

		handleError(w, err)
		return
	}

	resp.CountPerPage = gelbooru.PostsPerPage
	resp.TotalCount = results.TotalCount
	resp.Results = results.Posts
	respBody, err := json.Marshal(resp)
	if err != nil {
		handleError(w, err)
		return
	}

	if !isNaughty {
		writePostsToCache(query, page, respBody)
	}

	w.Write(respBody)
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
