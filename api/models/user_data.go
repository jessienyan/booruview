package models

import (
	"slices"
	"time"
)

const (
	SearchHistoryLimit = 100
)

type SearchHistoryEntry struct {
	SearchedAt time.Time `json:"searched_at"`
	Tags       []string  `json:"tags"`
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
		(*h)[0].SearchedAt = item.SearchedAt
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
	FavoritePosts []any         `json:"favorite_posts"`
	FavoriteTags  []any         `json:"favorite_tags"`
	Blacklist     []any         `json:"blacklist"`
	SearchHistory SearchHistory `json:"search_history"`
}
