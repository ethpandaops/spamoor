package handlers

import (
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

// PprofLogin renders the login gate for the /debug/pprof pages. The server
// serves it to a browser that asks for the profiling pages without a valid
// login; the page hands the session's token over as a short-lived cookie
// scoped to /debug/pprof and reloads, or asks for a login first.
func (fh *FrontendHandler) PprofLogin(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"debug/pprof_login.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "pprof", r.URL.Path, "Profiling", templateFiles)
	data.Data = &struct{}{}

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "pprof.go", "PprofLogin", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}
