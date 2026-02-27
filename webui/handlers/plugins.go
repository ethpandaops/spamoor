package handlers

import (
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

// Plugins renders the plugins management page.
func (fh *FrontendHandler) Plugins(w http.ResponseWriter, r *http.Request) {
	templateFiles := append(server.LayoutTemplateFiles,
		"plugins/plugins.html",
	)

	pageTemplate := server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "plugins", "/plugins", "Plugins", templateFiles)

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "plugins.go", "Plugins", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}
