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

func NewCDNHostResponse() CDNHostResponse {
	if api.UseMediaProxy {
		return CDNHostResponse{
			Image:      api.MediaProxyHost + "/?to=",
			Video:      api.MediaProxyHost + "/?to=",
			MediaProxy: true,
		}
	}

	hosts := gelbooru.GetCDNHosts()
	return CDNHostResponse{
		Image: hosts.ImageHost,
		Video: hosts.VideoHost,
	}
}

func CDNHostHandler(w http.ResponseWriter, req *http.Request) {
	respondJson(w, 200, NewCDNHostResponse())
}
