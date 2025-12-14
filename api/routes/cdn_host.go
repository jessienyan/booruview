package routes

import (
	"net/http"

	"codeberg.org/jessienyan/booruview/gelbooru"
)

type CDNHostResponse struct {
	Image string `json:"image"`
	Video string `json:"video"`
}

func CDNHostHandler(w http.ResponseWriter, req *http.Request) {
	hosts := gelbooru.GetCDNHosts()
	respondJson(w, 200, CDNHostResponse{
		Image: hosts.ImageHost,
		Video: hosts.VideoHost,
	})
}
