package gelbooru

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"

	api "github.com/jessienyan/booruview"
)

type fakeClient struct{}

func NewFakeClient() Client {
	return fakeClient{}
}

func (f fakeClient) ListPosts(tags string, page int) (*PostList, error) {
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
		results.Posts = append(results.Posts, newFakePost(r))
	}

	return &results, nil
}

func (f fakeClient) ListTags(tags string) ([]api.TagResponse, error) {
	panic("not impl")
	// seed := hash(tags)
}

func (f fakeClient) SearchTags(query string) ([]api.TagResponse, error) {
	panic("not impl")
	// seed := hash(query)
}

// hash uses the djb2 algorithm
// http://www.cse.yorku.ca/~oz/hash.html
func hash(s string) int64 {
	hash := int64(5381)
	for _, b := range []byte(s) {
		hash = ((hash << 5) + hash) + int64(b)
	}
	return hash
}

func randChoice[T any](r *rand.Rand, items []T) T {
	return items[r.Intn(len(items))]
}

func randChoiceN[T any](r *rand.Rand, items []T, count int) []T {
	if count > len(items) {
		count = len(items)
	}

	itemsCopy := make([]T, len(items))
	copy(itemsCopy, items)
	r.Shuffle(len(itemsCopy), func(i, j int) {
		itemsCopy[i], itemsCopy[j] = itemsCopy[j], itemsCopy[i]
	})

	return itemsCopy[:count]
}

func randRange(r *rand.Rand, min int, max int) int {
	return r.Int()%(max-min) + min
}

func randString(r *rand.Rand, src string, len int) string {
	var result strings.Builder

	for i := 0; i < len; i++ {
		result.WriteByte(randChoice(r, []byte(src)))
	}

	return result.String()
}

var (
	postRatings = []string{
		"explicit",
		"questionable",
		"sensitive",
		"general",
	}

	sampleTags = []string{
		":d",
		"1girl",
		"2girls",
		"absurdres",
		"ahoge",
		"ambience_synesthesia",
		"animal_ears",
		"apron",
		"bikini",
		"black_cape",
		"black_dress",
		"black_gloves",
		"black_hat",
		"black_pantyhose",
		"blue_apple",
		"blue_dress",
		"blue_eyes",
		"blunt_bangs",
		"blush",
		"bow_(music)",
		"bow",
		"breasts",
		"brown_hair",
		"candle",
		"cape",
		"chinese_commentary",
		"column",
		"commentary_request",
		"detached_sleeves",
		"dress",
		"ears_through_headwear",
		"fang",
		"film_grain",
		"frilled_apron",
		"frills",
		"garter_straps",
		"ghost_girl",
		"gloves",
		"gradient_hair",
		"green_dress",
		"green_hair",
		"green_sleeves",
		"hair_intakes",
		"hair_ornament",
		"hat",
		"heart_hands",
		"heart",
		"high_heels",
		"highres",
		"holding_bow_(music)",
		"holding_instrument",
		"holding_violin",
		"holding",
		"instrument",
		"jewelry",
		"long_hair",
		"looking_at_viewer",
		"looking_to_the_side",
		"maid_headdress",
		"maid",
		"multicolored_hair",
		"multiple_girls",
		"multiple_rings",
		"music",
		"official_alternate_costume",
		"on_back",
		"open_mouth",
		"pantyhose",
		"picture_frame",
		"pillar",
		"playing_instrument",
		"puffy_short_sleeves",
		"puffy_sleeves",
		"rabbit_ears",
		"ring",
		"short_sleeves",
		"sidelocks",
		"simple_background",
		"single_garter_strap",
		"single_glove",
		"single_thighhigh",
		"sleeveless_dress",
		"sleeveless",
		"smile",
		"solo",
		"swimsuit",
		"thigh_strap",
		"thighhighs",
		"triangular_headpiece",
		"two_side_up",
		"two-sided_cape",
		"two-sided_fabric",
		"violin",
		"white_apron",
		"white_background",
		"white_cape",
		"white_hair",
		"white_headdress",
		"witch_hat",
	}

	alphaNumericLower = "abcdefghijklmnopqrstuvwxyz0123456789"
)

const (
	maxThumbSize  = 250
	maxLowresSize = 850
)

func newFakePost(r *rand.Rand) api.PostResponse {
	var hashInput [8]byte
	r.Read(hashInput[:])
	hashData := md5.Sum(hashInput[:])
	digest := hex.EncodeToString(hashData[:])

	var thumbW, thumbH, lowresW, lowresH int
	width := randRange(r, 800, 3000)
	height := randRange(r, 800, 3000)
	aspect := float64(width) / float64(height)

	if width < height {
		thumbW = int(float64(maxThumbSize) * aspect)
		thumbH = maxThumbSize
		lowresW = int(float64(maxLowresSize) * aspect)
		lowresH = maxLowresSize
	} else {
		thumbW = maxThumbSize
		thumbH = int(float64(maxThumbSize) / aspect)
		lowresW = maxLowresSize
		lowresH = int(float64(maxLowresSize) / aspect)
	}

	resp := api.PostResponse{
		Id:                 randRange(r, 5_000_000, 13_000_000),
		CreatedAtTimestamp: time.Now().Unix() - int64(randRange(r, 0, 1_000_000)),
		Score:              randRange(r, 0, 50),
		Rating:             randChoice(r, postRatings),
		SourceUrl: fmt.Sprintf(
			"https://twitter.com/%s/status/%d",
			randString(r, alphaNumericLower, randRange(r, 5, 15)),
			randRange(r, 1930000000000000000, 1970000000000000000),
		),
		Uploader:        "danbooru",
		UploaderUrl:     "https://gelbooru.com/index.php?page=account&s=profile&id=6498",
		Tags:            randChoiceN(r, sampleTags, randRange(r, 5, 40)),
		ThumbnailUrl:    fmt.Sprintf("https://img4.gelbooru.com/thumbnails/%s/%s/thumbnail_%s.jpg", digest[:2], digest[2:4], digest),
		ThumbnailWidth:  thumbW,
		ThumbnailHeight: thumbH,
		LowResUrl:       fmt.Sprintf("https://img4.gelbooru.com/samples/%s/%s/sample_%s.jpg", digest[:2], digest[2:4], digest),
		LowResWidth:     lowresW,
		LowResHeight:    lowresH,
		ImageUrl:        fmt.Sprintf("https://img4.gelbooru.com/images/%s/%s/%s.jpg", digest[:2], digest[2:4], digest),
		Width:           width,
		Height:          height,
	}

	slices.Sort(resp.Tags)
	resp.Tags = append(resp.Tags, fmt.Sprintf("rating:%s", resp.Rating))

	return resp
}
