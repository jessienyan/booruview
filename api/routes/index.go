package routes

import (
	"html/template"
	"net/http"
	"os"

	api "codeberg.org/jessienyan/booruview"
	"codeberg.org/jessienyan/booruview/models"
	"github.com/rs/zerolog/log"
)

const (
	indexTemplatePath = "/dist/index.html"
)

var indexTemplate *template.Template

func IndexHandler(w http.ResponseWriter, req *http.Request) {
	readTemplateFile := indexTemplate == nil || api.DevMode
	if readTemplateFile {
		f, err := os.ReadFile(indexTemplatePath)
		if err != nil {
			log.Err(err).Msg("failed to read index.html")
			respondWithInternalError(w, err)
			return
		}

		indexTemplate, err = template.New("index").Parse(string(f))
		if err != nil {
			log.Err(err).Msg("failed to parse template")
			respondWithInternalError(w, err)
			return
		}
	}

	type TemplateContext struct {
		AccountData *models.UserDataJSON
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
		tmplContext.AccountData = &data
	}

	w.Header().Add("Content-Type", "text/html")
	indexTemplate.Execute(w, tmplContext)
}
