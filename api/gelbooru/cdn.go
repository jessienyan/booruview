package gelbooru

import (
	"net/url"
	"sync"
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

func UpdateCDNHosts() error {
	c := NewClient()

	images, err := c.ListPosts("-video", 1)
	if err != nil {
		return err
	}

	videos, err := c.ListPosts("video", 1)
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

	return nil
}
