package routes

import (
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/posts", PostsHandler)
	r.HandleFunc("/tags", TagsHandler)
	r.HandleFunc("/tagsearch", TagSearchHandler)
	r.HandleFunc("/settings/import", SettingImportHandler).Methods("POST")
	r.HandleFunc("/settings/export", SettingExportHandler).Methods("POST")
	r.HandleFunc("/hosts", CDNHostHandler)
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		type versionResponse struct {
			Version string `json:"version"`
		}
		respondJson(w, http.StatusOK, versionResponse{Version: api.AppVersion})
	})
	r.Use(EmptyResponseMiddleware, RecoverMiddleware)

	return r
}
