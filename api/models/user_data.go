package models

import (
	"encoding/json"
	"slices"
	"time"
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

type FavoritePost struct {
	CreatedAt time.Time `json:"created_at"`
	Post      any       `json:"post"`
}

type FavoriteTag struct {
	CreatedAt time.Time `json:"created_at"`
	Tag       any       `json:"tag"`
}

type BlacklistEntry struct {
	CreatedAt time.Time `json:"created_at"`
	Tag       any       `json:"tag"`
}

type UserDataJSON struct {
	FavoritePosts []FavoritePost   `json:"favorite_posts"`
	FavoriteTags  []FavoriteTag    `json:"favorite_tags"`
	Blacklist     []BlacklistEntry `json:"blacklist"`
	SearchHistory SearchHistory    `json:"search_history"`
}

func (ud UserDataJSON) MarshalJSON() ([]byte, error) {
	if ud.Blacklist == nil {
		ud.Blacklist = []BlacklistEntry{}
	}
	if ud.FavoritePosts == nil {
		ud.FavoritePosts = []FavoritePost{}
	}
	if ud.FavoriteTags == nil {
		ud.FavoriteTags = []FavoriteTag{}
	}
	if ud.SearchHistory == nil {
		ud.SearchHistory = SearchHistory{}
	}

	// Use a different type for marshalling, otherwise this will go into an infinite loop
	type marshalType UserDataJSON
	return json.Marshal(marshalType(ud))
}
