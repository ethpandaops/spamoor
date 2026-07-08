package spamoor

import (
	"fmt"
)

// WalletInfo describes one spamoor-managed wallet with its pool context. It is
// the shared enumeration format used by the daemon wallet overview and by
// cross-cutting scenarios (e.g. the ens-names wallet naming service).
type WalletInfo struct {
	Wallet *Wallet
	// Scenario is the scenario name the wallet belongs to ("" for the root
	// wallet, or when the host process does not know it - CLI mode).
	Scenario string
	// SpammerID is the daemon spammer id (0 in CLI mode and for the root wallet).
	SpammerID uint64
	// Name is the wallet name within its pool: "root", a child wallet index
	// ("1".."n") or a well-known wallet name.
	Name string
}

// WalletInfoProvider enumerates all wallets known to the host process. The
// daemon registers one on every wallet pool it creates; it stays nil in CLI
// mode.
type WalletInfoProvider func() []*WalletInfo

// SetWalletInfoProvider registers the host-wide wallet enumerator used by
// GetAllWalletInfos.
func (pool *WalletPool) SetWalletInfoProvider(provider WalletInfoProvider) {
	pool.walletInfoProvider = provider
}

// GetWalletInfos returns this pool's wallets (child and well-known, without
// the root wallet) with their pool-local names and the pool's spammer id.
func (pool *WalletPool) GetWalletInfos() []*WalletInfo {
	infos := make([]*WalletInfo, 0, len(pool.childWallets)+len(pool.wellKnownNames))
	for idx, wallet := range pool.childWallets {
		infos = append(infos, &WalletInfo{
			Wallet:    wallet,
			SpammerID: pool.GetSpammerID(),
			Name:      fmt.Sprintf("%d", idx+1),
		})
	}

	for _, config := range pool.wellKnownNames {
		wallet := pool.wellKnownWallets[config.Name]
		if wallet == nil {
			continue
		}

		infos = append(infos, &WalletInfo{
			Wallet:    wallet,
			SpammerID: pool.GetSpammerID(),
			Name:      config.Name,
		})
	}

	return infos
}

// GetAllWalletInfos returns all wallets known to the host process. In daemon
// mode this delegates to the registered provider (the root wallet plus every
// running spammer's pool); without a provider (CLI mode) it falls back to the
// root wallet plus this pool's own wallets.
func (pool *WalletPool) GetAllWalletInfos() []*WalletInfo {
	if pool.walletInfoProvider != nil {
		return pool.walletInfoProvider()
	}

	infos := make([]*WalletInfo, 0, len(pool.childWallets)+len(pool.wellKnownNames)+1)
	infos = append(infos, &WalletInfo{
		Wallet: pool.rootWallet.GetWallet(),
		Name:   "root",
	})

	return append(infos, pool.GetWalletInfos()...)
}
