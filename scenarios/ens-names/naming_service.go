package ensnames

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// namingRetryLimit gives up naming an address after this many reverted
// attempts (e.g. persistent races with another naming service instance).
const namingRetryLimit = 3

// namingWalletDuration is the registration duration of naming service names.
const namingWalletDuration = 365 * 24 * 3600

// namingCandidate is one wallet queued for naming in the current tick.
type namingCandidate struct {
	addr  common.Address
	label string
}

// runNamingService is the opt-out background routine that registers a fully
// wired ENS name (forward + addr.reverse resolution) for every spamoor wallet:
// the root wallet, this scenario's wallets and - in "all" mode - every other
// running spammer's wallets, discovered via the host wallet registry. All
// naming txs are sent by the scenario's own "nameservice" wallet through the
// permissionless SpamRegistrarController, so no transaction is ever sent from
// (and no nonce consumed on) the wallets being named.
func (s *Scenario) runNamingService(ctx context.Context) {
	nameSvcWallet := s.walletPool.GetWellKnownWallet("nameservice")
	if nameSvcWallet == nil {
		s.logger.Errorf("naming service: nameservice wallet not available")
		return
	}

	seen := make(map[common.Address]bool, 64)
	retries := make(map[common.Address]int, 16)
	namedCount := 0

	ticker := time.NewTicker(scenario.GlobalSlotDuration)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			named, err := s.namingServiceTick(ctx, nameSvcWallet, seen, retries)
			if err != nil {
				s.logger.Warnf("naming service: %v", err)
				continue
			}
			if named > 0 {
				namedCount += named
				s.logger.Infof("naming service: named %d wallets (%d total)", named, namedCount)
			}
		}
	}
}

// namingServiceTick collects unnamed wallets and names up to NamingPerSlot of
// them, returning the number of successfully named wallets.
func (s *Scenario) namingServiceTick(ctx context.Context, nameSvcWallet *spamoor.Wallet, seen map[common.Address]bool, retries map[common.Address]int) (int, error) {
	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return 0, scenario.ErrNoClients
	}

	candidates, err := s.collectNamingCandidates(ctx, seen, retries)
	if err != nil {
		return 0, err
	}
	if len(candidates) == 0 {
		return 0, nil
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)
	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return 0, err
	}

	txs := make([]*types.Transaction, 0, len(candidates))
	for _, candidate := range candidates {
		addr := candidate.addr
		label := candidate.label

		// Gas is estimated per tx: ENS registrations create lots of fresh state
		// and the chain's EIP-8037 cost-per-state-byte moves over time, so any
		// static budget would eventually run out of gas.
		tx, err := nameSvcWallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
			GasFeeCap: uint256.MustFromBig(feeCap),
			GasTipCap: uint256.MustFromBig(tipCap),
			Value:     uint256.NewInt(0),
		}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
			return s.deployment.SpamController.RegisterNamed(transactOpts, label, addr, big.NewInt(namingWalletDuration))
		})
		if err != nil {
			return 0, fmt.Errorf("could not build registerNamed tx for %s: %w", label, err)
		}

		txs = append(txs, tx)
	}

	receipts, err := s.walletPool.GetTxPool().SendTransactionBatch(ctx, nameSvcWallet, txs, &spamoor.BatchOptions{
		SendTransactionOptions: spamoor.SendTransactionOptions{
			Client:      client,
			ClientGroup: s.options.ClientGroup,
			Rebroadcast: true,
		},
		MaxRetries:   2,
		PendingLimit: 10,
	})
	if err != nil {
		return 0, fmt.Errorf("could not send naming txs: %w", err)
	}

	named := 0
	for i, receipt := range receipts {
		addr := candidates[i].addr
		if receipt != nil && receipt.Status == types.ReceiptStatusSuccessful {
			seen[addr] = true
			named++
			s.logger.Debugf("naming service: named %s as %s.eth", addr.Hex(), candidates[i].label)
			continue
		}

		// reverted or lost: usually a race with another naming service
		// instance; the reverse-record pre-check resolves it next tick
		retries[addr]++
		if retries[addr] >= namingRetryLimit {
			s.logger.Warnf("naming service: giving up on %s (%s.eth) after %d attempts", addr.Hex(), candidates[i].label, retries[addr])
			seen[addr] = true
		}
	}

	return named, nil
}

// collectNamingCandidates returns up to NamingPerSlot wallets that still need
// a name, checking chain state so restarts (and concurrent naming services)
// never re-name a wallet.
func (s *Scenario) collectNamingCandidates(ctx context.Context, seen map[common.Address]bool, retries map[common.Address]int) ([]*namingCandidate, error) {
	callOpts := &bind.CallOpts{Context: ctx}

	candidates := make([]*namingCandidate, 0, s.options.NamingPerSlot)
	for _, info := range s.collectWalletInfos() {
		if uint64(len(candidates)) >= s.options.NamingPerSlot {
			break
		}

		addr := info.Wallet.GetAddress()
		if seen[addr] {
			continue
		}

		// already has an addr.reverse record (from an earlier run or another
		// naming service instance)?
		resolver, err := s.deployment.Registry.Resolver(callOpts, addrReverseNode2LD(addr))
		if err != nil {
			return nil, fmt.Errorf("could not check reverse record of %s: %w", addr.Hex(), err)
		}
		if resolver != (common.Address{}) {
			seen[addr] = true
			continue
		}

		label := s.walletEnsLabel(info)

		// label taken but reverse record missing: registered for a different
		// address in an earlier incarnation, cannot be re-registered
		available, err := s.deployment.Base.Available(callOpts, new(big.Int).SetBytes(crypto.Keccak256([]byte(label))))
		if err != nil {
			return nil, fmt.Errorf("could not check availability of %s: %w", label, err)
		}
		if !available {
			s.logger.Warnf("naming service: label %s.eth already taken by another address, skipping %s", label, addr.Hex())
			seen[addr] = true
			continue
		}

		if retries[addr] == 0 {
			s.logger.Debugf("naming service: queueing %s as %s.eth", addr.Hex(), label)
		}
		candidates = append(candidates, &namingCandidate{addr: addr, label: label})
	}

	return candidates, nil
}

// collectWalletInfos enumerates the wallets in scope for the naming service:
// every wallet the host process knows about in "all" mode, or just the root
// wallet plus this scenario's own pool in "pool" mode.
func (s *Scenario) collectWalletInfos() []*spamoor.WalletInfo {
	if s.options.WalletNaming == walletNamingAll {
		return s.walletPool.GetAllWalletInfos()
	}

	infos := make([]*spamoor.WalletInfo, 0, s.walletPool.GetWalletCount()+2)
	infos = append(infos, &spamoor.WalletInfo{
		Wallet: s.walletPool.GetRootWallet().GetWallet(),
		Name:   "root",
	})

	return append(infos, s.walletPool.GetWalletInfos()...)
}

// walletEnsLabel derives the ENS label for a wallet:
// root wallet         -> spamoor-root
// daemon mode wallets -> spamoor-<scenario><spammerid>-<name/index>
// CLI mode wallets    -> spamoor-<scenario>-<name/index>
func (s *Scenario) walletEnsLabel(info *spamoor.WalletInfo) string {
	if info.Name == "root" && info.Scenario == "" {
		return "spamoor-root"
	}

	scenarioName := info.Scenario
	if scenarioName == "" {
		scenarioName = ScenarioName
	}

	if info.SpammerID != 0 {
		return fmt.Sprintf("spamoor-%s%d-%s", scenarioName, info.SpammerID, info.Name)
	}

	return fmt.Sprintf("spamoor-%s-%s", scenarioName, info.Name)
}

// addrReverseNode2LD returns the registry node of <addr>.addr.reverse (the
// label is the lowercase hex address without 0x prefix, per ENSIP-3).
func addrReverseNode2LD(addr common.Address) [32]byte {
	label := hex.EncodeToString(addr.Bytes())
	labelHash := crypto.Keccak256([]byte(label))

	node := [32]byte{}
	copy(node[:], crypto.Keccak256(addrReverseNode[:], labelHash))
	return node
}
