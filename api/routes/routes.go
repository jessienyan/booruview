package routes

import (
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	client := gelbooru.NewClient()
	r := mux.NewRouter()

	r.Use(RecoverMiddleware)

	r.Handle("/tags", TagsHandler{Client: client})
	r.Handle("/tagsearch", TagSearchHandler{Client: client})
	r.HandleFunc("/settings/import", SettingImportHandler).Methods("POST")
	r.HandleFunc("/settings/export", SettingExportHandler).Methods("POST")
	r.Handle("/posts", PostsHandler{Client: client})
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
	authRouter.Use(AuthMiddleware)
	authRouter.HandleFunc("/account", AccountHandler).Methods("GET", "PATCH", "DELETE")
	authRouter.HandleFunc("/account/password", ChangePasswordHandler).Methods("POST")

	return r
}
