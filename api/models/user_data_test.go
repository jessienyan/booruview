package models_test

import (
	"testing"
	"time"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"codeberg.org/jessienyan/booruview/testutil"
	"github.com/stretchr/testify/require"
)

func TestSearchHistoryList_Clean_SortsByNewestFirst(t *testing.T) {
	now := testutil.Time()
	oldTime := now.Add(-time.Hour)
	list := models.SearchHistoryList{
		{Date: oldTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
	}

	list.Clean()

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
		{Date: oldTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Clean_RemovesDuplicateQueries(t *testing.T) {
	now := testutil.Time()
	oldTime := now.Add(-time.Hour)
	list := models.SearchHistoryList{
		{Date: oldTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}

	list.Clean()

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Clean_SkipsEmptyQueries(t *testing.T) {
	now := testutil.Time()
	list := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{}},
	}

	list.Clean()

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Remove_ByQuery(t *testing.T) {
	now := testutil.Time()
	list := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
	}

	list.Remove([]string{"a"})

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Remove_EmptyQueryList(t *testing.T) {
	now := testutil.Time()
	list := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}

	list.Remove([]string{})

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Add_DeduplicatesAndSorts(t *testing.T) {
	now := testutil.Time()
	oldTime := now.Add(-2 * time.Hour)
	midTime := now.Add(-time.Hour)

	list := models.SearchHistoryList{
		{Date: oldTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: midTime, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
	}

	newEntries := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "c"}}}},
	}

	list.Add(newEntries)

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "c"}}}},
		{Date: midTime, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
	}
	require.Equal(t, expected, list)
}

func TestSearchHistoryList_Remove_DeduplicatesAndSorts(t *testing.T) {
	now := testutil.Time()
	oldTime := now.Add(-2 * time.Hour)
	midTime := now.Add(-time.Hour)

	list := models.SearchHistoryList{
		{Date: oldTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: midTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
		{Date: midTime, Query: models.SearchQuery{Include: api.TagList{{Name: "b"}}}},
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "c"}}}},
	}

	list.Remove([]string{"b"})

	expected := models.SearchHistoryList{
		{Date: now, Query: models.SearchQuery{Include: api.TagList{{Name: "c"}}}},
		{Date: midTime, Query: models.SearchQuery{Include: api.TagList{{Name: "a"}}}},
	}
	require.Equal(t, expected, list)
}
