package models

import (
	"encoding/json"
	"slices"
	"time"

	api "codeberg.org/jessienyan/booruview"
)

const (
	SearchHistoryLimit = 100
)

type SearchHistoryEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
}

// Items in the search history are sorted by newest first.
type SearchHistory []SearchHistoryEntry

// Adds a new search entry. To keep the history tidy, every entry must be unique.
// When adding an entry that already exists, the existing entry is removed.
func (h *SearchHistory) Add(item SearchHistoryEntry) {
	slices.Sort(item.Tags)

	i := slices.IndexFunc(*h, func(cmp SearchHistoryEntry) bool {
		return slices.Equal(item.Tags, cmp.Tags)
	})

	// If repeating the last search, just update the timestamp
	if i == 0 {
		(*h)[0].CreatedAt = item.CreatedAt
		return
	}

	newHistory := make(SearchHistory, 0, SearchHistoryLimit)
	newHistory = append(newHistory, item)

	if i == -1 {
		newHistory = append(newHistory, *h...)
	} else {
		// Append the existing history *except* the duplicate entry. We're basically
		// popping the entry and moving it to the front of the list
		newHistory = append(newHistory, (*h)[:i]...)
		newHistory = append(newHistory, (*h)[i+1:]...)
	}

	*h = newHistory[:min(len(newHistory), SearchHistoryLimit)]
}

type UserDataJSON struct {
	FavoritePosts []api.PostResponse `json:"favorite_posts" validate:"dive"`
	FavoriteTags  []api.TagResponse  `json:"favorite_tags" validate:"dive"`
	Blacklist     []api.TagResponse  `json:"blacklist" validate:"dive"`
	SearchHistory SearchHistory      `json:"search_history"`
}

func (ud UserDataJSON) MarshalJSON() ([]byte, error) {
	if ud.Blacklist == nil {
		ud.Blacklist = []api.TagResponse{}
	}
	if ud.FavoritePosts == nil {
		ud.FavoritePosts = []api.PostResponse{}
	}
	if ud.FavoriteTags == nil {
		ud.FavoriteTags = []api.TagResponse{}
	}
	if ud.SearchHistory == nil {
		ud.SearchHistory = SearchHistory{}
	}

	// Use a different type for marshalling, otherwise this will go into an infinite loop
	type marshalType UserDataJSON
	return json.Marshal(marshalType(ud))
}
