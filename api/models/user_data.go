package models

import (
	"encoding/json"
	"time"

	api "codeberg.org/jessienyan/booruview"
)

const (
	SearchHistoryLimit = 100
)

type SearchHistoryEntry struct {
	Date    time.Time         `json:"date" validate:"required"`
	Include []api.TagResponse `json:"include"`
	Exclude []api.TagResponse `json:"exclude"`
}

type UserDataJSON struct {
	FavoritePosts []api.PostResponse   `json:"favorite_posts" validate:"dive"`
	FavoriteTags  []api.TagResponse    `json:"favorite_tags" validate:"dive"`
	Blacklist     []api.TagResponse    `json:"blacklist" validate:"dive"`
	SearchHistory []SearchHistoryEntry `json:"search_history" validate:"dive"`
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
		ud.SearchHistory = []SearchHistoryEntry{}
	}

	// Use a different type for marshalling, otherwise this will go into an infinite loop
	type marshalType UserDataJSON
	return json.Marshal(marshalType(ud))
}
