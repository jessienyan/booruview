package api_test

import (
	"testing"

	api "codeberg.org/jessienyan/booruview"
	"github.com/stretchr/testify/require"
)

func TestPostList_Add_MergesWithoutDuplicates(t *testing.T) {
	list := api.PostList{{Id: 1}, {Id: 2}}
	list.Add(api.PostList{{Id: 2}, {Id: 3}})

	expected := api.PostList{{Id: 2}, {Id: 3}, {Id: 1}}
	require.Equal(t, expected, list)
}

func TestPostList_Add_PlacesNewPostsAtBeginning(t *testing.T) {
	list := api.PostList{{Id: 1}}
	list.Add(api.PostList{{Id: 2}})

	expected := api.PostList{{Id: 2}, {Id: 1}}
	require.Equal(t, expected, list)
}

func TestPostList_Remove_ByID(t *testing.T) {
	list := api.PostList{{Id: 1}, {Id: 2}, {Id: 3}}
	list.Remove([]int{2})

	expected := api.PostList{{Id: 1}, {Id: 3}}
	require.Equal(t, expected, list)
}

func TestPostList_Remove_Empty(t *testing.T) {
	list := api.PostList{{Id: 1}, {Id: 2}}
	list.Remove([]int{})

	expected := api.PostList{{Id: 1}, {Id: 2}}
	require.Equal(t, expected, list)
}

func TestPostList_Clean_DeduplicatesByID(t *testing.T) {
	list := api.PostList{{Id: 1}, {Id: 2}, {Id: 1}}
	list.Clean()

	expected := api.PostList{{Id: 1}, {Id: 2}}
	require.Equal(t, expected, list)
}

func TestTagList_Add_MergesWithoutDuplicates(t *testing.T) {
	list := api.TagList{{Name: "a"}, {Name: "b"}}
	list.Add(api.TagList{{Name: "b"}, {Name: "c"}})

	expected := api.TagList{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	require.Equal(t, expected, list)
}

func TestTagList_Clean_SortsAlphabetically(t *testing.T) {
	list := api.TagList{{Name: "c"}, {Name: "a"}, {Name: "b"}}
	list.Clean()

	expected := api.TagList{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	require.Equal(t, expected, list)
}

func TestTagList_Clean_RemovesDuplicateNames(t *testing.T) {
	list := api.TagList{{Name: "a"}, {Name: "b"}, {Name: "a"}}
	list.Clean()

	expected := api.TagList{{Name: "a"}, {Name: "b"}}
	require.Equal(t, expected, list)
}

func TestTagList_Remove_ByName(t *testing.T) {
	list := api.TagList{{Name: "a"}, {Name: "b"}, {Name: "c"}}
	list.Remove([]string{"b"})

	expected := api.TagList{{Name: "a"}, {Name: "c"}}
	require.Equal(t, expected, list)
}

func TestTagList_Remove_Empty(t *testing.T) {
	list := api.TagList{{Name: "a"}, {Name: "b"}}
	list.Remove([]string{})

	expected := api.TagList{{Name: "a"}, {Name: "b"}}
	require.Equal(t, expected, list)
}

func TestTagList_Equal_CompareByName(t *testing.T) {
	list1 := api.TagList{{Name: "a", Count: 1}, {Name: "b", Count: 2}}
	list2 := api.TagList{{Name: "a", Count: 10}, {Name: "b", Count: 20}}
	require.True(t, list1.Equal(list2))
}
