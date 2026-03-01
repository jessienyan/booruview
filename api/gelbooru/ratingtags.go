package gelbooru

import (
	"context"

	api "codeberg.org/jessienyan/booruview"
)

var (
	ratingTags = api.TagList{
		{Name: "rating:general", Type: api.Metadata},
		{Name: "rating:sensitive", Type: api.Metadata},
		{Name: "rating:questionable", Type: api.Metadata},
		{Name: "rating:explicit", Type: api.Metadata},
	}
)

// AddRatingTagsToValkey inserts some overrides for rating tags into the cache.
// Gelbooru doesn't recognize these as actual tags
func AddRatingTagsToValkey() {
	vk := api.Valkey()

	for _, tag := range ratingTags {
		key := TagCacheKey(tag.Name)

		vk.Do(context.Background(),
			vk.B().
				Set(). // Don't use a TTL since this is meant as an override
				Key(key).
				Value(api.TagToCacheValue(tag)).
				Build())
	}
}
