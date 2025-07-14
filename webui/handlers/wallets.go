package handlers

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethpandaops/spamoor/webui/server"
)

type WalletsPage struct {
	RootWallet     *WalletInfo   `json:"root_wallet"`
	SpammerWallets []*WalletInfo `json:"spammer_wallets"`
}

type WalletInfo struct {
	Address        string `json:"address"`
	Name           string `json:"name"`
	Balance        string `json:"balance"`
	PendingNonce   uint64 `json:"pending_nonce"`
	ConfirmedNonce uint64 `json:"confirmed_nonce"`
	SpammerID      int64  `json:"spammer_id,omitempty"`
	SpammerName    string `json:"spammer_name,omitempty"`
	SpammerStatus  int    `json:"spammer_status,omitempty"`
}

func (fh *FrontendHandler) Wallets(w http.ResponseWriter, r *http.Request) {
	var templateFiles = append(server.LayoutTemplateFiles,
		"wallets/wallets.html",
	)

	var pageTemplate = server.GetTemplate(templateFiles...)
	data := server.InitPageData(r, "wallets", "/wallets", "Spamoor Wallets", templateFiles)

	var pageError error
	data.Data, pageError = fh.getWalletsPageData()
	if pageError != nil {
		server.HandlePageError(w, r, pageError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if server.HandleTemplateError(w, r, "wallets.go", "Wallets", "", pageTemplate.ExecuteTemplate(w, "layout", data)) != nil {
		return
	}
}

func (fh *FrontendHandler) getWalletsPageData() (*WalletsPage, error) {
	// Get root wallet info
	rootWallet := fh.daemon.GetRootWallet()
	if rootWallet == nil {
		return nil, fmt.Errorf("root wallet not found")
	}

	rootAddr := rootWallet.GetWallet().GetAddress()
	rootBalance := rootWallet.GetWallet().GetBalance()
	rootPendingNonce := rootWallet.GetWallet().GetNonce()
	rootConfirmedNonce := rootWallet.GetWallet().GetConfirmedNonce()

	data := &WalletsPage{
		RootWallet: &WalletInfo{
			Address:        rootAddr.String(),
			Name:           "root",
			Balance:        weiToEthString(rootBalance),
			PendingNonce:   rootPendingNonce,
			ConfirmedNonce: rootConfirmedNonce,
		},
		SpammerWallets: []*WalletInfo{},
	}

	// Get all spammers (not just running ones)
	spammers := fh.daemon.GetAllSpammers()
	for _, spammer := range spammers {
		walletPool := spammer.GetWalletPool()
		if walletPool == nil {
			continue
		}

		for _, wallet := range walletPool.GetAllWallets() {
			addr := wallet.GetAddress()
			balance := wallet.GetBalance()
			pendingNonce := wallet.GetNonce()
			confirmedNonce := wallet.GetConfirmedNonce()

			data.SpammerWallets = append(data.SpammerWallets, &WalletInfo{
				Address:        addr.String(),
				Name:           walletPool.GetWalletName(addr),
				Balance:        weiToEthString(balance),
				PendingNonce:   pendingNonce,
				ConfirmedNonce: confirmedNonce,
				SpammerID:      spammer.GetID(),
				SpammerName:    spammer.GetName(),
				SpammerStatus:  spammer.GetStatus(),
			})
		}
	}

	return data, nil
}

func weiToEthString(wei *big.Int) string {
	if wei == nil {
		return "0"
	}
	// Convert wei to eth (1 eth = 10^18 wei)
	fbalance := new(big.Float).SetInt(wei)
	ethValue := new(big.Float).Quo(fbalance, new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)))
	// Format with fixed decimal places to avoid scientific notation
	return ethValue.Text('f', 6)
}
