package routes

import (
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(RecoverMiddleware, LoadUserMiddleware)
	r.HandleFunc("/tags", TagsHandler)
	r.HandleFunc("/tagsearch", TagSearchHandler)
	r.HandleFunc("/settings/import", SettingImportHandler).Methods("POST")
	r.HandleFunc("/settings/export", SettingExportHandler).Methods("POST")
	r.HandleFunc("/hosts", CDNHostHandler)
	r.HandleFunc("/version", func(w http.ResponseWriter, req *http.Request) {
		type versionResponse struct {
			Version string `json:"version"`
		}
		respondJson(w, http.StatusOK, versionResponse{Version: api.AppVersion})
	})
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(LoadUserMiddleware)
	authRouter.HandleFunc("/account", AccountHandler)
	authRouter.HandleFunc("/posts", PostsHandler)

	return r
}
