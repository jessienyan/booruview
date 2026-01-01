package routes

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/pkg/errors"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/valkey-io/valkey-go"
)

type TagSearchResponse struct {
	Results []api.TagResponse `json:"results"`
}

func TagSearchHandler(w http.ResponseWriter, req *http.Request) {
	if isRateLimited(w, req, tagSearchApiCost) {
		return
	}

	resp := TagSearchResponse{
		Results: []api.TagResponse{},
	}

	query := cleanTag(req.FormValue("q"))
	query = strings.ReplaceAll(query, " ", "_")

	if query == "" {
		respondWithBadRequest(w, "required GET param `q` is missing or blank")
		return
	}

	cached, err := getCachedTagSearch(query)
	if err != nil {
		respondWithInternalError(w, err)
		return
	}

	// Cache hit
	if cached != nil {
		respondJson(w, 200, api.DecompressData(cached))
		return
	}

	results, err := gelbooru.NewClient().SearchTags(query)
	if err != nil {
		if errors.As(err, &gelbooru.GelbooruError{}) {
			respondWithGelbooruUnavailable(w)
			return
		}

		respondWithInternalError(w, err)
		return
	}

	resp.Results = results
	respData := respondJson(w, http.StatusOK, resp)
	writeTagSearchToCache(query, respData)
}

func getCachedTagSearch(query string) ([]byte, error) {
	vk := api.Valkey()
	cached := vk.Do(context.Background(),
		vk.B().
			Get().
			Key(gelbooru.TagSearchCacheKey(query)).
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

func writeTagSearchToCache(query string, data []byte) error {
	vk := api.Valkey()
	buf := bytes.Buffer{}
	if err := api.CompressData(&buf, data); err != nil {
		return err
	}

	return vk.Do(context.Background(),
		vk.B().
			Setex().
			Key(gelbooru.TagSearchCacheKey(query)).
			Seconds(api.TagSearchTtl).
			Value(buf.String()).
			Build(),
	).Error()
}
