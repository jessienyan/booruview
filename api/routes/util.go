package routes

import (
	"slices"
	"strings"
)

// Returns a normalized version of the tag with
func cleanTag(tag string) string {
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

	if len(tag) == 0 || tag == "-" {
		return ""
	}

	tag = strings.ToLower(tag)
	tag = strings.ReplaceAll(tag, " ", "_")

	return tag
}

// Returns a new []string of normalized and sorted tags
func cleanTagList(tags []string) []string {
	cleaned := make([]string, 0, len(tags))

	for _, t := range tags {
		t = cleanTag(t)

		if t != "" {
			cleaned = append(cleaned, t)
		}
	}

	slices.Sort(cleaned)
	return slices.Compact(cleaned)
}
