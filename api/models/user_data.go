package models

import "time"

type SearchHistoryEntry struct {
	SearchedAt time.Time `json:"searched_at"`
	Tags       []string  `json:"tags"`
}

type UserDataJSON struct {
	FavoritePosts []any                `json:"favorite_posts"`
	FavoriteTags  []any                `json:"favorite_tags"`
	Blacklist     []any                `json:"blacklist"`
	SearchHistory []SearchHistoryEntry `json:"search_history"`
}
