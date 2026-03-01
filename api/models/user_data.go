package models

import (
	"encoding/json"
	"slices"
	"strings"
	"time"

	api "codeberg.org/jessienyan/booruview"
)

const (
	SearchHistoryLimit = 100
)

type SearchHistoryEntry struct {
	Date  time.Time `json:"date" validate:"required"`
	Query struct {
		Include api.TagList `json:"include" validate:"required"`
		Exclude api.TagList `json:"exclude" validate:"required"`
	} `json:"query"`
}

// Clean normalizes the query tags: removing duplicates and sorting them
func (entry *SearchHistoryEntry) Clean() {
}

// Tags returns
func (entry SearchHistoryEntry) Tags() string {
	tags := strings.Builder{}
	for _, t := range entry.Query.Include {
		tags.WriteString(t.Name)
		tags.WriteByte(',')
	}
	for _, t := range entry.Query.Exclude {
		tags.WriteString(t.Name)
		tags.WriteByte(',')
	}
	return tags.String()
}

type UserDataJSON struct {
	FavoritePosts []api.PostResponse   `json:"favorite_posts" validate:"dive"`
	FavoriteTags  api.TagList          `json:"favorite_tags" validate:"dive"`
	Blacklist     api.TagList          `json:"blacklist" validate:"dive"`
	SearchHistory []SearchHistoryEntry `json:"search_history" validate:"dive"`
}

func (ud UserDataJSON) MarshalJSON() ([]byte, error) {
	if ud.Blacklist == nil {
		ud.Blacklist = api.TagList{}
	}
	if ud.FavoritePosts == nil {
		ud.FavoritePosts = []api.PostResponse{}
	}
	if ud.FavoriteTags == nil {
		ud.FavoriteTags = api.TagList{}
	}
	if ud.SearchHistory == nil {
		ud.SearchHistory = []SearchHistoryEntry{}
	}

	// Use a different type for marshalling, otherwise this will go into an infinite loop
	type marshalType UserDataJSON
	return json.Marshal(marshalType(ud))
}

// Removes duplicate entries
func (ud *UserDataJSON) Clean() {
	postIds := make(map[int]struct{}, len(ud.FavoritePosts))
	ud.FavoritePosts = slices.DeleteFunc(ud.FavoritePosts, func(p api.PostResponse) bool {
		_, dupe := postIds[p.Id]
		if dupe {
			return true
		}
		postIds[p.Id] = struct{}{}
		return false
	})

	tagNames := make(map[string]struct{}, len(ud.FavoriteTags))
	ud.FavoriteTags = slices.DeleteFunc(ud.FavoriteTags, func(p api.TagResponse) bool {
		_, dupe := tagNames[p.Name]
		if dupe {
			return true
		}
		tagNames[p.Name] = struct{}{}
		return false
	})

	tagNames = make(map[string]struct{}, len(ud.Blacklist))
	ud.Blacklist = slices.DeleteFunc(ud.Blacklist, func(p api.TagResponse) bool {
		_, dupe := tagNames[p.Name]
		if dupe {
			return true
		}
		tagNames[p.Name] = struct{}{}
		return false
	})

	searchTags := make(map[string]struct{}, len(ud.SearchHistory))
	ud.SearchHistory = slices.DeleteFunc(ud.SearchHistory, func(h SearchHistoryEntry) bool {
		tags := h.Tags()
		_, dupe := searchTags[tags]
		if dupe {
			return true
		}
		searchTags[tags] = struct{}{}
		return false
	})
}
