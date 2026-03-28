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

type SearchQueryNames struct {
	Include []string `json:"include" validate:"required"`
	Exclude []string `json:"exclude" validate:"required"`
}

func (query *SearchQueryNames) Clean() {
	slices.Sort(query.Include)
	query.Include = slices.Compact(query.Include)
	slices.Sort(query.Exclude)
	query.Include = slices.Compact(query.Exclude)
}

func (query SearchQueryNames) Tags() string {
	tags := strings.Builder{}
	for _, t := range query.Include {
		tags.WriteString(t)
		tags.WriteByte(',')
	}
	for _, t := range query.Exclude {
		tags.WriteByte('-')
		tags.WriteString(t)
		tags.WriteByte(',')
	}
	return tags.String()
}

// Clean normalizes the query tags to be sorted and de-duped
func (query *SearchQuery) Clean() {
	query.Include.Clean()
	query.Exclude.Clean()
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

// Clean sorts and de-dupes the search history
func (lst *SearchHistoryList) Clean() {
	for _, entry := range *lst {
		entry.Clean()
	}

	// Sort by newest first
	slices.SortFunc(*lst, func(a, b SearchHistoryEntry) int {
		return b.Date.Compare(a.Date)
	})

	// Remove entries that have duplicate queries, leaving only the most
	// recent entry
	queries := make(map[string]struct{}, len(*lst))
	*lst = slices.DeleteFunc(*lst, func(entry SearchHistoryEntry) bool {
		tags := entry.Query.Tags()
		_, seen := queries[tags]

		if seen {
			return true
		}

		queries[tags] = struct{}{}
		return false
	})
}

// queries should be the result of calling .Tags() on a tag list
func (lst *SearchHistoryList) Remove(queries []string) {
	if len(queries) == 0 {
		return
	}

	lookup := make(map[string]struct{}, len(queries))
	for _, q := range queries {
		lookup[q] = struct{}{}
	}

	*lst = slices.DeleteFunc(*lst, func(entry SearchHistoryEntry) bool {
		_, shouldDelete := lookup[entry.Query.Tags()]
		return shouldDelete
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

func (ud *UserDataJSON) Clean() {
	ud.FavoritePosts.Clean()
	ud.FavoriteTags.Clean()
	ud.Blacklist.Clean()
	ud.SearchHistory.Clean()
}
