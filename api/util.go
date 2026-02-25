package api

import (
	"slices"
	"strings"
)

// Returns a normalized version of the tag
func CleanTag(tag string) string {
	hyphens := 0
	for _, c := range tag {
		if c != '-' {
			break
		}
		hyphens++
	}

	// Trim excess hyphens
	if hyphens > 1 {
		tag = "-" + tag[hyphens:]
	}

	tag = strings.TrimSpace(tag)

	if tag == "-" {
		return ""
	}

	tag = strings.ToLower(tag)

	// Not an OR tag
	if tag[0] != '{' {
		tag = strings.ReplaceAll(tag, " ", "_")
	}

	return tag
}

// Returns a new []string of normalized and sorted tags
func CleanTagList(tags []string) []string {
	cleaned := make([]string, 0, len(tags))

	for _, t := range tags {
		t = CleanTag(t)

		if t != "" {
			cleaned = append(cleaned, t)
		}
	}

	slices.Sort(cleaned)
	return slices.Compact(cleaned)
}

// Sorts and dedupes the tag list, returning the new slice
func CleanTagResponseList(tags []TagResponse) []TagResponse {
	slices.SortFunc(tags, func(a, b TagResponse) int {
		return strings.Compare(a.Name, b.Name)
	})
	tags = slices.CompactFunc(tags, func(a, b TagResponse) bool {
		return a.Name == b.Name
	})
	return tags
}
