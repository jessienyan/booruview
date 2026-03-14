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

type SearchQuery struct {
	Include api.TagList `json:"include" validate:"required"`
	Exclude api.TagList `json:"exclude" validate:"required"`
}

// Clean normalizes the query tags to be sorted and de-duped
func (query *SearchQuery) Clean() {
	query.Include = query.Include.Clean()
	query.Exclude = query.Exclude.Clean()
}

// Equal returns whether the two queries are identical.
// Both queries should be cleaned before calling Equal
func (query SearchQuery) Equal(other SearchQuery) bool {
	return query.Include.Equal(other.Include) && query.Exclude.Equal(other.Exclude)
}

// Tags return a normalized list of comma-separated tags.
// The query should be cleaned before calling Tags
func (query SearchQuery) Tags() string {
	tags := strings.Builder{}
	for _, t := range query.Include {
		tags.WriteString(t.Name)
		tags.WriteByte(',')
	}
	for _, t := range query.Exclude {
		tags.WriteByte('-')
		tags.WriteString(t.Name)
		tags.WriteByte(',')
	}
	return tags.String()
}

type SearchHistoryEntry struct {
	Date  time.Time   `json:"date" validate:"required"`
	Query SearchQuery `json:"query"`
}

// Clean normalizes the query tags to be sorted and de-duped
func (entry *SearchHistoryEntry) Clean() {
	entry.Query.Clean()
}

type SearchHistoryList []SearchHistoryEntry

func (lst *SearchHistoryList) Remove(queries []string) {
	if len(queries) == 0 {
		return
	}

	*lst = slices.DeleteFunc(*lst, func(entry SearchHistoryEntry) bool {
		for _, q := range queries {
			if entry.Query.Tags() == q {
				return true
			}
		}
		return false
	})
}

type UserDataJSON struct {
	FavoritePosts api.PostList      `json:"favorite_posts" validate:"dive"`
	FavoriteTags  api.TagList       `json:"favorite_tags" validate:"dive"`
	Blacklist     api.TagList       `json:"blacklist" validate:"dive"`
	SearchHistory SearchHistoryList `json:"search_history" validate:"dive"`
}

func (ud UserDataJSON) MarshalJSON() ([]byte, error) {
	if ud.Blacklist == nil {
		ud.Blacklist = api.TagList{}
	}
	if ud.FavoritePosts == nil {
		ud.FavoritePosts = api.PostList{}
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
