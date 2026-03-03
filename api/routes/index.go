package routes

import (
	"html/template"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

var indexTemplate *template.Template

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	if indexTemplate == nil {
		f, err := os.ReadFile("../index.html")
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read index.html")
		}

		indexTemplate, err = template.New("index").Parse(string(f))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse template")
		}
	}

	var templateContext any
	user := getUser(req)
	if user == nil {
		templateContext = nil
	} else {
		data, err := user.Data.ParseJSON()
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		templateContext = data
	}

	w.Header().Add("Content-Type", "text/html")
	indexTemplate.Execute(w, templateContext)
}
