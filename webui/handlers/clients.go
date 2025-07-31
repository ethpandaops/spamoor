package handlers

import (
	"context"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ethpandaops/spamoor/webui/server"
)

type ClientsPage struct {
	Clients     []*ClientsPageClient `json:"clients"`
	ClientCount uint64               `json:"client_count"`
}

type ClientsPageClient struct {
	Index         int      `json:"index"`
	Name          string   `json:"name"`
	Group         string   `json:"group"`  // First group for backward compatibility
	Groups        []string `json:"groups"` // All groups
	Version       string   `json:"version"`
	BlockHeight   uint64   `json:"block_height"`
	IsReady       bool     `json:"ready"`
	Enabled       bool     `json:"enabled"`
	NameOverride  string   `json:"name_override,omitempty"`
	TotalRequests uint64   `json:"total_requests"`
	TxRequests    uint64   `json:"tx_requests"`
	RpcFailures   uint64   `json:"rpc_failures"`
}

// Clients will return the "clients" page using a go template
func (fh *FrontendHandler) Clients(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"clients/clients.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "clients", "/clients", "Clients", templateFiles)

	var pageError error
	data.Data, pageError = fh.getClientsPageData(r.Context())
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "clients.go", "Clients", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return // an error has occurred and was processed
	}
}

func (fh *FrontendHandler) getClientsPageData(ctx context.Context) (*ClientsPage, error) {
	pageData := &ClientsPage{
		Clients: []*ClientsPageClient{},
	}

	goodClients := fh.daemon.GetClientPool().GetAllGoodClients()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	wg := sync.WaitGroup{}

	// Get all clients from pool
	for idx, client := range fh.daemon.GetClientPool().GetAllClients() {
		blockHeight, _ := client.GetLastBlockHeight()
		totalReqs, txReqs, rpcFails := client.GetRequestStats()

		clientData := &ClientsPageClient{
			Index:         idx,
			Name:          client.GetName(),
			Group:         client.GetClientGroup(),
			Groups:        client.GetClientGroups(),
			BlockHeight:   blockHeight,
			IsReady:       slices.Contains(goodClients, client),
			Enabled:       client.IsEnabled(),
			NameOverride:  client.GetNameOverride(),
			TotalRequests: totalReqs,
			TxRequests:    txReqs,
			RpcFailures:   rpcFails,
		}

		wg.Add(1)
		go func(clientData *ClientsPageClient) {
			version, err := client.GetClientVersion(ctx)
			if err != nil {
				version = "Unknown"
			}

			clientData.Version = version
			wg.Done()
		}(clientData)

		pageData.Clients = append(pageData.Clients, clientData)
	}
	pageData.ClientCount = uint64(len(pageData.Clients))

	wg.Wait()

	return pageData, nil
}
