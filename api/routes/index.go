package routes

import (
	"html/template"
	"net/http"
	"os"

	"codeberg.org/jessienyan/booruview/models"
	"github.com/rs/zerolog/log"
)

const (
	indexTemplatePath = "/index.html"
)

var indexTemplate *template.Template

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	if indexTemplate == nil {
		f, err := os.ReadFile(indexTemplatePath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to read index.html")
		}

		indexTemplate, err = template.New("index").Parse(string(f))
		if err != nil {
			log.Fatal().Err(err).Msg("failed to parse template")
		}
	}

	type TemplateContext struct {
		AccountData models.UserDataJSON
		CDNHosts    CDNHostResponse
	}
	tmplContext := TemplateContext{
		CDNHosts: NewCDNHostResponse(),
	}

	user := GetUser(req)
	if user != nil {
		data, err := user.Data.ParseJSON()
		if err != nil {
			respondWithInternalError(w, err)
			return
		}
		tmplContext.AccountData = data
	}

	w.Header().Add("Content-Type", "text/html")
	indexTemplate.Execute(w, tmplContext)
}
