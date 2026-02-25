package routes

import (
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
)

type CDNHostResponse struct {
	Image      string `json:"image"`
	Video      string `json:"video"`
	MediaProxy bool   `json:"media_proxy"`
}

func CDNHostHandler(w http.ResponseWriter, req *http.Request) {
	if api.UseMediaProxy {
		respondJson(w, 200, CDNHostResponse{
			Image:      api.MediaProxyHost + "/?to=",
			Video:      api.MediaProxyHost + "/?to=",
			MediaProxy: true,
		})
		return
	}

	hosts := gelbooru.GetCDNHosts()
	respondJson(w, 200, CDNHostResponse{
		Image: hosts.ImageHost,
		Video: hosts.VideoHost,
	})
}
