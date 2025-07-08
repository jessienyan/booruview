package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	api "github.com/jessienyan/booruview"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()
	r.Path("/posts").HandlerFunc(PostsHandler)
	r.Path("/tags").HandlerFunc(TagsHandler)
	r.Path("/tagsearch").HandlerFunc(TagSearchHandler)
	r.Path("/version").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("{\"version\": \"%s\"}", api.AppVersion)))
	})
	r.Use(RecoverMiddleware, RateLimitMiddleware)

	return r
}
