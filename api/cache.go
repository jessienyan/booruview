package api

import (
	"fmt"
	"strconv"
	"strings"
)

func TagToCacheValue(tag TagResponse) string {
	return fmt.Sprintf("%s,%d", tag.Type, tag.Count)
}

func TagFromCacheValue(tagName string, val string) (tag TagResponse, err error) {
	parts := strings.Split(val, ",")
	if len(parts) != 2 {
		err = fmt.Errorf("TagFromCacheValue: expected value to have 2 fields (has %d)", len(parts))
		return
	}

	var count int
	count, err = strconv.Atoi(parts[1])
	if err != nil {
		return
	}

	tag = TagResponse{
		Name:  tagName,
		Type:  ParseTagType(parts[0]),
		Count: count,
	}

	return
}
