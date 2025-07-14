package handlers

import (
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

type GraphsPage struct {
	// Add any server-side data needed for the graphs page
}

// Graphs will return the "graphs" page using a go template
func (fh *FrontendHandler) Graphs(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"graphs/graphs.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "graphs", "/graphs", "Transaction Graphs", templateFiles)

	var pageError error
	data.Data, pageError = fh.getGraphsPageData()
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "graphs.go", "Graphs", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getGraphsPageData() (*GraphsPage, error) {
	pageData := &GraphsPage{
		// Initialize any server-side data if needed
	}

	return pageData, nil
}
