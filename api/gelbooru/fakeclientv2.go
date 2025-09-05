package gelbooru

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"

	api "codeberg.org/jessienyan/booruview"
)

var (
	recentPosts     = make([]api.PostResponse, 0, 1000)
	recentPostIndex = 0
)

func CollectRecentPosts(posts []api.PostResponse) {
	if len(posts) == 0 {
		return
	}

	// Just sample half
	half := posts[:len(posts)/2]
	for _, p := range half {
		// Append to recentPosts until it's full. Once it's full, switch to using an index
		// to treat it like a circular array
		if len(recentPosts) < cap(recentPosts) {
			recentPosts = append(recentPosts, p)
		} else {
			recentPosts[recentPostIndex] = p
			recentPostIndex = (recentPostIndex + 1) % len(recentPosts)
		}

	}
}

type fakeClientv2 struct{}

// v2 client returns random posts that have been recently searched. The v1 fake client
// often returns images that are broken or deleted
func NewFakeClientv2() Client {
	return fakeClientv2{}
}

func (f fakeClientv2) ListPosts(tags string, page int) (*PostList, error) {
	var r *rand.Rand

	// Generate total count from tags
	countSeed := hash(tags)
	r = rand.New(rand.NewSource(countSeed))
	count := randRange(r, 100, 10_000)

	resultSeed := hash(fmt.Sprintf("%s,%d", tags, page))
	results := PostList{
		TotalCount: count,
		Posts:      make([]api.PostResponse, 0, PostsPerPage),
	}

	r.Seed(resultSeed)

	for i := 0; i < PostsPerPage; i++ {
		randomPost := recentPosts[r.Int()%len(recentPosts)]
		results.Posts = append(results.Posts, randomPost)
	}

	// Sort the results by time
	slices.SortFunc(results.Posts, func(a, b api.PostResponse) int {
		return cmp.Compare(a.CreatedAtTimestamp, b.CreatedAtTimestamp)
	})
	slices.Reverse(results.Posts)

	return &results, nil
}

func (f fakeClientv2) ListTags(tags string) ([]api.TagResponse, error) {
	panic("not impl")
}

func (f fakeClientv2) SearchTags(query string) ([]api.TagResponse, error) {
	panic("not impl")
}
