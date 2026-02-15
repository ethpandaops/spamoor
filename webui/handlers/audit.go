package handlers

import (
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

// Audit will return the "audit" page using a go template
// The page content is loaded dynamically via JavaScript and protected API calls
func (fh *FrontendHandler) Audit(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"audit/audit.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "audit", "/audit", "Audit Logs", templateFiles)

	// No server-side data - page loads content dynamically via protected API
	data.Data = &struct{}{}

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "audit.go", "Audit", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}
