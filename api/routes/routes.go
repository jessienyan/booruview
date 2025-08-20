package routes

import (
	"fmt"
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
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"version\": \"%s\"}", api.AppVersion)))
	})
	r.Use(RecoverMiddleware)

	return r
}
