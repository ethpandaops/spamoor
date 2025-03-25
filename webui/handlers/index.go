package handlers

import (
	"net/http"
	"time"

	"github.com/ethpandaops/spamoor/webui"
)

type IndexPage struct {
	Spammers []*IndexPageSpammer
}

type IndexPageSpammer struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Scenario    string    `json:"scenario"`
	Status      int       `json:"status"`
	Config      string    `json:"config"`
	CreatedAt   time.Time `json:"created_at"`
}

// Index will return the "index" page using a go template
func (fh *FrontendHandler) Index(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(webui.LayoutTemplateFiles,
		"index/index.html",
	)

	var pageTemplate = webui.GetTemplate(templateFiles...)
	data := webui.InitPageData(w, r, "index", "/", "Spamoor Dashboard", templateFiles)

	var pageError error
	data.Data, pageError = fh.getIndexPageData()
	if pageError != nil {
		webui.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if webui.HandleTemplateError(w, r, "index.go", "Index", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getIndexPageData() (*IndexPage, error) {
	spammers := fh.daemon.GetAllSpammers()
	models := make([]*IndexPageSpammer, len(spammers))

	for i, s := range spammers {
		models[i] = &IndexPageSpammer{
			ID:          s.GetID(),
			Name:        s.GetName(),
			Description: s.GetDescription(),
			Scenario:    s.GetScenario(),
			Status:      s.GetStatus(),
			Config:      s.GetConfig(),
			CreatedAt:   time.Unix(s.GetCreatedAt(), 0),
		}
	}

	return &IndexPage{
		Spammers: models,
	}, nil
}
