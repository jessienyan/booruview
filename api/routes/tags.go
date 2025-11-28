package routes

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/valkey-io/valkey-go"
)

const (
	tagLimit = 100
)

type TagsResponse struct {
	Results []api.TagResponse `json:"results"`
}

func TagsHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, tagApiCost) {
		return
	}

	if err := req.ParseForm(); err != nil {
		respondWithInternalError(w, err)
		return
	}

	query := cleanTagList(req.Form["t"])

	// Strip leading hyphen
	for i := range query {
		if query[i][0] == '-' {
			query[i] = query[i][1:]
		}
	}

	// write empty response
	if len(query) == 0 {
		resp := TagsResponse{Results: []api.TagResponse{}}
		respondJson(w, http.StatusOK, resp)
		return
	}

	if len(query) > tagLimit {
		respondWithBadRequest(w, fmt.Sprintf("limit of %d tags", tagLimit))
		return
	}

	cached, cachedMap, err := getCachedTags(query)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	var missing []string
	for _, query := range query {
		if _, ok := cachedMap[query]; !ok {
			missing = append(missing, query)
		}
	}

	var tags []api.TagResponse
	if len(missing) > 0 {
		tags, err = gelbooru.NewClient().ListTags(strings.Join(missing, " "))
		if err != nil {
			if errors.As(err, &gelbooru.GelbooruError{}) {
				respondWithGelbooruUnavailable(w)
				return
			}

			respondWithInternalError(w, err)
			return
		}
	}

	resp := TagsResponse{Results: append(cached, tags...)}
	respondJson(w, http.StatusOK, resp)

	if len(tags) > 0 {
		writeCachedTags(tags)
	}
}

func writeCachedTags(tags []api.TagResponse) {
	vk := api.Valkey()
	cmds := make(valkey.Commands, 0, len(tags))

	for _, tag := range tags {
		key := gelbooru.TagCacheKey(tag.Name)

		cmds = append(cmds,
			vk.B().
				Setex().
				Key(key).
				Seconds(api.TagTtl).
				Value(api.TagToCacheValue(tag)).
				Build())
	}
	vk.DoMulti(context.Background(), cmds...)
}

func getCachedTags(query []string) ([]api.TagResponse, map[string]api.TagResponse, error) {
	keys := make([]string, len(query))
	for i, query := range query {
		keys[i] = gelbooru.TagCacheKey(query)
	}

	vk := api.Valkey()
	cached := vk.Do(context.Background(),
		vk.B().
			Mget().
			Key(keys...).
			Build(),
	)

	if err := cached.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, nil, nil
		} else {
			return nil, nil, err
		}
	}

	entries, err := cached.AsStrSlice()
	if err != nil {
		return nil, nil, err
	}

	resp := make([]api.TagResponse, 0, len(entries))
	respMap := make(map[string]api.TagResponse, len(entries))
	for i, entry := range entries {
		if entry == "" {
			continue
		}

		tag, err := api.TagFromCacheValue(query[i], entry)
		if err != nil {
			log.Err(err).Msg("failed to parse tag from cache")
			continue
		}

		resp = append(resp, tag)
		respMap[tag.Name] = tag
	}

	return resp, respMap, nil
}
