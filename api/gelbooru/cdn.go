package gelbooru

import (
	"net/url"
	"sync"

	api "codeberg.org/jessienyan/booruview"
	"github.com/rs/zerolog/log"
)

type GelbooruCDN struct {
	ImageHost string
	VideoHost string
}

var (
	cdn      GelbooruCDN
	cdnMutex sync.RWMutex
)

func GetCDNHosts() GelbooruCDN {
	cdnMutex.RLock()
	defer cdnMutex.RUnlock()

	return cdn
}

func UpdateCDNHosts(client GelbooruClient) error {
	if api.UseMediaProxy {
		return nil
	}

	images, err := client.ListPosts("-video", 1)
	if err != nil {
		return err
	}

	videos, err := client.ListPosts("video", 1)
	if err != nil {
		return err
	}

	imgUrl, _ := url.Parse(images.Posts[0].ImageUrl)
	videoUrl, _ := url.Parse(videos.Posts[0].ImageUrl)

	cdnMutex.Lock()
	defer cdnMutex.Unlock()

	cdn = GelbooruCDN{
		ImageHost: imgUrl.Hostname(),
		VideoHost: videoUrl.Hostname(),
	}

	log.Info().Str("img", cdn.ImageHost).Str("video", cdn.VideoHost).Msg("refreshed gelbooru CDN hosts")

	return nil
}
