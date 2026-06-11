package routes

import (
	"net/http"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/gelbooru"
	"github.com/gorilla/mux"
)

func NewRouter(client gelbooru.GelbooruClient) *mux.Router {
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

	maybeAuthRouter := r.NewRoute().Subrouter()
	maybeAuthRouter.Use(AuthMiddleware)
	maybeAuthRouter.HandleFunc("/index.html", IndexHandler)

	authRouter := r.NewRoute().Subrouter()
	authRouter.Use(AuthMiddleware, RequireAuthMiddleware)
	authRouter.HandleFunc("/account", AccountDeleteHandler).Methods("DELETE")
	authRouter.HandleFunc("/account/data", AccountDataGetHandler).Methods("GET")
	authRouter.HandleFunc("/account/data", AccountDataPatchHandler).Methods("PATCH")
	authRouter.HandleFunc("/account/data", AccountDataPutHandler).Methods("PUT")
	authRouter.HandleFunc("/account/password", ChangePasswordHandler).Methods("POST")

	return r
}
