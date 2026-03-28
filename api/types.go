package api

import (
	"slices"
	"strings"
)

type TagType string

const (
	Tag        TagType = "tag"
	Artist     TagType = "artist"
	Copyright  TagType = "copyright"
	Character  TagType = "character"
	Metadata   TagType = "metadata"
	Deprecated TagType = "deprecated"
	Unknown    TagType = "unknown"
)

func ParseTagType(val string) TagType {
	switch val {
	case string(Tag):
		return Tag
	case string(Artist):
		return Artist
	case string(Copyright):
		return Copyright
	case string(Character):
		return Character
	case string(Metadata):
		return Metadata
	case string(Deprecated):
		return Deprecated
	default:
		return Unknown
	}
}

type TagResponse struct {
	Name  string  `json:"name" validate:"required"`
	Type  TagType `json:"type" validate:"required"`
	Count int     `json:"count"`
}

type TagList []TagResponse

// Clean modifies the TagList by sorting and removing duplicates
func (lst *TagList) Clean() {
	lst.Sort()
	*lst = slices.CompactFunc(*lst, func(a, b TagResponse) bool {
		return a.Name == b.Name
	})
}

func (lst TagList) Equal(other TagList) bool {
	return slices.EqualFunc(lst, other, func(a, b TagResponse) bool {
		return a.Name == b.Name
	})
}

// Sort sorts the TagList in-place.
func (lst TagList) Sort() {
	slices.SortFunc(lst, func(a, b TagResponse) int {
		return strings.Compare(a.Name, b.Name)
	})
}

// Remove modifies the TagList by removing tags that match by name
func (lst *TagList) Remove(tagNames []string) {
	if len(tagNames) == 0 {
		return
	}

	lookup := make(map[string]struct{}, len(tagNames))
	for _, t := range tagNames {
		lookup[t] = struct{}{}
	}

	*lst = slices.DeleteFunc(*lst, func(t TagResponse) bool {
		_, shouldDelete := lookup[t.Name]
		return shouldDelete
	})

	lst.Clean()
}

type PostResponse struct {
	Id                 int      `json:"id"`
	CreatedAtTimestamp int64    `json:"created_at"`
	Score              int      `json:"score"`
	Rating             string   `json:"rating"`
	SourceUrl          string   `json:"source_url"`
	Uploader           string   `json:"uploader"`
	UploaderUrl        string   `json:"uploader_url"`
	Tags               []string `json:"tags" validate:"required"`
	ThumbnailUrl       string   `json:"thumbnail_url"`
	ThumbnailWidth     int      `json:"thumbnail_width"`
	ThumbnailHeight    int      `json:"thumbnail_height"`
	LowResUrl          string   `json:"lowres_url"`
	LowResWidth        int      `json:"lowres_width"`
	LowResHeight       int      `json:"lowres_height"`
	ImageUrl           string   `json:"image_url"`
	Width              int      `json:"width"`
	Height             int      `json:"height"`
}

// Clean sorts and removes dupes from the post's tag list
func (p PostResponse) Clean() PostResponse {
	slices.Sort(p.Tags)
	p.Tags = slices.Compact(p.Tags)
	return p
}

type PostList []PostResponse

// Clean removes duplicate posts
func (lst PostList) Clean() PostList {
	// Keep track of what post IDs we've seen as we go through the list. If we encounter a post ID
	// that was seen previously, drop it from the list
	idsSeen := make(map[int]struct{}, len(lst))
	lst = slices.DeleteFunc(lst, func(p PostResponse) bool {
		_, alreadySeen := idsSeen[p.Id]
		if alreadySeen {
			return true
		}
		idsSeen[p.Id] = struct{}{}
		return false
	})
	return lst
}

// Remove modifies the PostList by removing posts that match by id
func (lst *PostList) Remove(ids []int) {
	if len(ids) == 0 {
		return
	}

	lookup := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		lookup[id] = struct{}{}
	}

	*lst = slices.DeleteFunc(*lst, func(t PostResponse) bool {
		_, shouldDelete := lookup[t.Id]
		return shouldDelete
	})

	lst.Clean()
}
