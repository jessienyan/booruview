package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	api "codeberg.org/jessienyan/booruview"
	"github.com/stretchr/testify/require"
)

func init() {
	api.LoadEnv()
	if err := api.InitValkey(); err != nil {
		panic(err)
	}
}

func TestTagSearchHandler_EmptyQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/tagsearch?q=", nil)
	rec := httptest.NewRecorder()
	TagSearchHandler{}.ServeHTTP(rec, req)
	require.Equal(t, rec.Code, http.StatusBadRequest)
}
