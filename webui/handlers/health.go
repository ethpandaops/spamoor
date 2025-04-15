package handlers

import (
	"context"
	"net/http"
	"slices"

	"github.com/ethpandaops/spamoor/webui/server"
)

type HealthPage struct {
	Clients     []*HealthPageClient `json:"clients"`
	ClientCount uint64              `json:"client_count"`
}

type HealthPageClient struct {
	Index       int    `json:"index"`
	Name        string `json:"name"`
	Group       string `json:"group"`
	Version     string `json:"version"`
	BlockHeight uint64 `json:"block_height"`
	IsReady     bool   `json:"ready"`
}

// Health will return the "health" page using a go template
func (fh *FrontendHandler) Health(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"health/health.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "health", "/health", "Health", templateFiles)

	var pageError error
	data.Data, pageError = fh.getHealthPageData(r.Context())
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "health.go", "Health", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}

func (fh *FrontendHandler) getHealthPageData(ctx context.Context) (*HealthPage, error) {
	pageData := &HealthPage{
		Clients: []*HealthPageClient{},
	}

	goodClients := fh.daemon.GetClientPool().GetAllGoodClients()

	// Get all clients from pool
	for idx, client := range fh.daemon.GetClientPool().GetAllClients() {
		blockHeight, _ := client.GetLastBlockHeight()

		version, err := client.GetClientVersion(ctx)
		if err != nil {
			version = "Unknown"
		}

		clientData := &HealthPageClient{
			Index:       idx,
			Name:        client.GetName(),
			Group:       client.GetClientGroup(),
			Version:     version,
			BlockHeight: blockHeight,
			IsReady:     slices.Contains(goodClients, client),
		}

		pageData.Clients = append(pageData.Clients, clientData)
	}
	pageData.ClientCount = uint64(len(pageData.Clients))

	return pageData, nil
}
