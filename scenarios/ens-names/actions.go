package ensnames

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/scenarios/ens-names/contract"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// registerRetryLimit drops a committed name after this many failed register
// attempts (e.g. when the commitment expired at very low throughput).
const registerRetryLimit = 3

// rentPriceBufferPercent is added on top of the rentPrice() quote; overpayment
// is refunded by the controller.
const rentPriceBufferPercent = 5

// fillCommitChance is the percent chance a below-target wallet advances its
// registration pipeline instead of doing maintenance on the names it already
// owns. Below 100 so the fleet doesn't march in lockstep commit/register
// waves during the initial fill (wallets are hit round-robin, so with an
// unconditional pipeline-first rule every wallet commits in the same slot
// window and registers in the next one, yielding long single-action phases).
const fillCommitChance = 70

// churn registration duration bounds (seconds): short enough that names
// visibly expire during a run.
const (
	churnMinDuration = 60
	churnMaxDuration = 3600
)

// noopResult is the onResult callback for stateless actions.
func noopResult(bool) {}

// buildActionTx picks and builds the next ENS action for the given wallet:
// it advances the commit-reveal registration pipeline until the wallet owns
// NamesPerWallet names, then performs weighted-random maintenance operations
// on the owned names.
func (s *Scenario) buildActionTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, walletIdx uint64, txIdx uint64, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	state := s.getWalletState(wallet.GetAddress())
	rng := rand.New(rand.NewSource(int64(txIdx)*0x9e3779b1 + int64(walletIdx) + 1))

	state.mu.Lock()
	defer state.mu.Unlock()

	// registration pipeline first: reveal a mature commitment, or place a new
	// commitment while the wallet is below its name target. Once at target,
	// the rotation op in buildMaintenanceTx keeps feeding this pipeline.
	if state.pending != nil && !state.pending.inFlight {
		if !state.pending.committed {
			// commit tx failed and was not cleaned up (defensive)
			state.pending = nil
		} else if s.commitmentExpired(state.pending) {
			s.logger.Warnf("dropping name %s: commitment expired before registration", state.pending.label)
			state.pending = nil
		} else if s.commitmentMature(state.pending) {
			return s.buildRegisterTx(ctx, client, wallet, state, feeCap, tipCap)
		}
	}

	if state.pending == nil && uint64(len(state.names)) < s.options.NamesPerWallet {
		// with no names yet there is nothing to maintain; otherwise mix in
		// maintenance on the existing names to de-synchronize the fill waves
		if len(state.names) == 0 || rng.Intn(100) < fillCommitChance {
			return s.buildCommitTx(ctx, client, wallet, state, walletIdx, feeCap, tipCap)
		}
	}

	// reclaim received transfers before anything else so registry ownership
	// catches up with the ERC721 ownership.
	for _, name := range state.names {
		if name.needReclaim && !name.inFlight {
			return s.buildReclaimTx(ctx, client, wallet, state, name, feeCap, tipCap)
		}
	}

	return s.buildMaintenanceTx(ctx, client, wallet, state, walletIdx, rng, feeCap, tipCap)
}

// commitmentMature reports whether the commitment is old enough to reveal.
// A full slot of margin is added on top of the local confirmation time so
// the latest mined block (which the register gas estimation runs against)
// already considers the commitment mature.
func (s *Scenario) commitmentMature(name *nameState) bool {
	minAge := time.Duration(s.options.MinCommitmentAge) * time.Second
	return time.Since(name.commitTime) > minAge+scenario.GlobalSlotDuration+2*time.Second
}

// commitmentExpired reports whether the commitment aged past the controller's
// max commitment age (minus a slot of margin).
func (s *Scenario) commitmentExpired(name *nameState) bool {
	maxAge := time.Duration(s.options.MaxCommitmentAge) * time.Second
	return time.Since(name.commitTime) > maxAge-12*time.Second
}

// buildCommitTx starts the commit-reveal pipeline for a fresh name: it builds
// the v1.7.0 Registration struct, pre-checks availability, computes the
// commitment via the controller and submits commit().
func (s *Scenario) buildCommitTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, walletIdx uint64, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	label := s.newLabel(state, walletIdx)

	available, err := s.deployment.Controller.Available(&bind.CallOpts{Context: ctx}, label)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not check availability of %s: %w", label, err)
	}
	if !available {
		return nil, nil, "", fmt.Errorf("label %s unexpectedly unavailable", label)
	}

	secret := [32]byte{}
	if _, err := cryptorand.Read(secret[:]); err != nil {
		return nil, nil, "", fmt.Errorf("could not generate commitment secret: %w", err)
	}

	node := ethNameNode(label)
	setAddrData, err := packCall(contract.PublicResolverMetaData, "setAddr0", node, wallet.GetAddress())
	if err != nil {
		return nil, nil, "", err
	}

	// reverseRecord bit 2 sets the default (chain-agnostic) reverse record for
	// the wallet's first name. The addr.reverse record is deliberately left to
	// the wallet naming service, so the two never fight over it.
	reverseRecord := uint8(0)
	if len(state.names) == 0 && state.nameSeq == 1 {
		reverseRecord = 2
	}

	name := &nameState{
		label:   label,
		node:    node,
		tokenID: labelhash(label),
		registration: contract.IETHRegistrarControllerRegistration{
			Label:         label,
			Owner:         wallet.GetAddress(),
			Duration:      new(big.Int).SetUint64(s.options.RegistrationDuration),
			Secret:        secret,
			Resolver:      s.deployment.ResolverAddr,
			Data:          [][]byte{setAddrData},
			ReverseRecord: reverseRecord,
			Referrer:      [32]byte{},
		},
	}

	commitment, err := s.deployment.Controller.MakeCommitment(&bind.CallOpts{Context: ctx}, name.registration)
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not compute commitment for %s: %w", label, err)
	}

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Controller.Commit(transactOpts, commitment)
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build commit tx: %w", err)
	}

	name.inFlight = true
	state.pending = name

	onResult := func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		name.inFlight = false
		if success {
			name.committed = true
			name.commitTime = time.Now()
		} else {
			state.pending = nil
		}
	}

	return tx, onResult, fmt.Sprintf("commit %s.eth", label), nil
}

// buildRegisterTx reveals a mature commitment, paying the oracle rent price.
func (s *Scenario) buildRegisterTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	name := state.pending

	value, err := s.rentPriceWithBuffer(ctx, name.label, name.registration.Duration)
	if err != nil {
		return nil, nil, "", err
	}

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     value,
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Controller.Register(transactOpts, name.registration)
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build register tx: %w", err)
	}

	name.inFlight = true

	onResult := func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		name.inFlight = false
		if success {
			name.registered = true
			state.names = append(state.names, name)
			state.pending = nil
			s.retireOldestNames(state)
		} else {
			name.registerTry++
			if name.registerTry >= registerRetryLimit {
				s.logger.Warnf("dropping name %s after %d failed register attempts", name.label, name.registerTry)
				state.pending = nil
			}
		}
	}

	return tx, onResult, fmt.Sprintf("register %s.eth", name.label), nil
}

// retireOldestNames drops the oldest names from active management once the
// wallet exceeds its name target (rotation registrations push it above). The
// retired names stay registered on-chain, they just stop receiving
// maintenance operations. Must be called with state.mu held.
func (s *Scenario) retireOldestNames(state *walletState) {
	for uint64(len(state.names)) > s.options.NamesPerWallet {
		retired := false
		for i, name := range state.names {
			if name.inFlight {
				continue
			}

			state.names = append(state.names[:i], state.names[i+1:]...)
			s.logger.Debugf("retired %s.eth from active management", name.label)
			retired = true
			break
		}
		if !retired {
			// every name has an operation in flight; retry after the next
			// rotation registration
			return
		}
	}
}

// buildReclaimTx updates the registry ownership of a name received via ERC721
// transfer.
func (s *Scenario) buildReclaimTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Base.Reclaim(transactOpts, new(big.Int).SetBytes(name.tokenID[:]), wallet.GetAddress())
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build reclaim tx: %w", err)
	}

	name.inFlight = true

	onResult := func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		name.inFlight = false
		if success {
			name.needReclaim = false
		}
	}

	return tx, onResult, fmt.Sprintf("reclaim %s.eth", name.label), nil
}

// maintenanceOp describes one weighted maintenance action candidate.
type maintenanceOp struct {
	weight uint64
	build  func() (*types.Transaction, func(success bool), string, error)
}

// buildMaintenanceTx performs a weighted-random maintenance operation on the
// wallet's names. Operations that are infeasible for the current state (no
// unwrapped name to transfer, wrapper approval missing, ...) are excluded from
// the draw; short-lived churn registrations are the always-feasible fallback.
func (s *Scenario) buildMaintenanceTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, walletIdx uint64, rng *rand.Rand, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	pickName := func(filter func(*nameState) bool) *nameState {
		candidates := make([]*nameState, 0, len(state.names))
		for _, name := range state.names {
			if !name.inFlight && !name.needReclaim && (filter == nil || filter(name)) {
				candidates = append(candidates, name)
			}
		}
		if len(candidates) == 0 {
			return nil
		}
		return candidates[rng.Intn(len(candidates))]
	}

	unwrapped := func(n *nameState) bool { return !n.wrapped }
	wrapped := func(n *nameState) bool { return n.wrapped }

	ops := make([]maintenanceOp, 0, 8)

	// rotation: keep registering fresh names via the full commit-reveal flow
	// even when the wallet is at its name target; the oldest name is retired
	// from active management once the new one lands (opt-out with weight 0)
	if s.options.RotationWeight > 0 && state.pending == nil {
		ops = append(ops, maintenanceOp{s.options.RotationWeight, func() (*types.Transaction, func(success bool), string, error) {
			return s.buildCommitTx(ctx, client, wallet, state, walletIdx, feeCap, tipCap)
		}})
	}

	if s.options.RenewWeight > 0 {
		if name := pickName(nil); name != nil {
			ops = append(ops, maintenanceOp{s.options.RenewWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildRenewTx(ctx, client, wallet, state, name, feeCap, tipCap)
			}})
		}
	}

	if s.options.RecordUpdateWeight > 0 {
		if name := pickName(nil); name != nil {
			ops = append(ops, maintenanceOp{s.options.RecordUpdateWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildRecordUpdateTx(ctx, client, wallet, state, name, rng, feeCap, tipCap)
			}})
		}
	}

	if s.options.TransferWeight > 0 && s.walletPool.GetWalletCount() > 1 {
		if name := pickName(unwrapped); name != nil {
			ops = append(ops, maintenanceOp{s.options.TransferWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildTransferTx(ctx, client, wallet, state, name, walletIdx, rng, feeCap, tipCap)
			}})
		}
	}

	if s.options.AbandonWeight > 0 {
		if name := pickName(unwrapped); name != nil {
			ops = append(ops, maintenanceOp{s.options.AbandonWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildAbandonTx(ctx, client, wallet, state, name, feeCap, tipCap)
			}})
		}
	}

	if s.options.ReverseWeight > 0 {
		if name := pickName(nil); name != nil {
			ops = append(ops, maintenanceOp{s.options.ReverseWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildReverseUpdateTx(ctx, client, wallet, state, name, feeCap, tipCap)
			}})
		}
	}

	if s.options.WrapWeight > 0 {
		if !state.wrapperApproved && !state.approveInFlight {
			ops = append(ops, maintenanceOp{s.options.WrapWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildWrapperApprovalTx(ctx, client, wallet, state, feeCap, tipCap)
			}})
		} else if state.wrapperApproved {
			if name := pickName(unwrapped); name != nil {
				ops = append(ops, maintenanceOp{s.options.WrapWeight, func() (*types.Transaction, func(success bool), string, error) {
					return s.buildWrapTx(ctx, client, wallet, state, name, feeCap, tipCap)
				}})
			}
		}
		if name := pickName(wrapped); name != nil {
			ops = append(ops, maintenanceOp{s.options.WrapWeight, func() (*types.Transaction, func(success bool), string, error) {
				return s.buildUnwrapTx(ctx, client, wallet, state, name, feeCap, tipCap)
			}})
		}
	}

	if s.options.ChurnWeight > 0 || len(ops) == 0 {
		churnWeight := s.options.ChurnWeight
		if churnWeight == 0 {
			churnWeight = 1
		}
		ops = append(ops, maintenanceOp{churnWeight, func() (*types.Transaction, func(success bool), string, error) {
			return s.buildChurnTx(ctx, client, wallet, state, walletIdx, rng, feeCap, tipCap)
		}})
	}

	totalWeight := uint64(0)
	for _, op := range ops {
		totalWeight += op.weight
	}

	draw := rng.Uint64() % totalWeight
	for _, op := range ops {
		if draw < op.weight {
			return op.build()
		}
		draw -= op.weight
	}

	return ops[len(ops)-1].build()
}

// buildRenewTx extends a name's registration through the commit-reveal
// controller (paying rent).
func (s *Scenario) buildRenewTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	duration := new(big.Int).SetUint64(s.options.RenewalDuration)

	value, err := s.rentPriceWithBuffer(ctx, name.label, duration)
	if err != nil {
		return nil, nil, "", err
	}

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     value,
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Controller.Renew(transactOpts, name.label, duration, [32]byte{})
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build renew tx: %w", err)
	}

	return tx, s.nameResult(state, name, nil), fmt.Sprintf("renew %s.eth", name.label), nil
}

// buildRecordUpdateTx updates the addr record and a text record of a name via
// the resolver's multicall.
func (s *Scenario) buildRecordUpdateTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, rng *rand.Rand, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	textKeys := []string{"url", "avatar", "description", "com.github"}
	key := textKeys[rng.Intn(len(textKeys))]

	setAddrData, err := packCall(contract.PublicResolverMetaData, "setAddr0", name.node, wallet.GetAddress())
	if err != nil {
		return nil, nil, "", err
	}
	setTextData, err := packCall(contract.PublicResolverMetaData, "setText", name.node, key, fmt.Sprintf("spamoor-%d", rng.Uint32()))
	if err != nil {
		return nil, nil, "", err
	}

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Resolver.MulticallWithNodeCheck(transactOpts, name.node, [][]byte{setAddrData, setTextData})
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build record update tx: %w", err)
	}

	return tx, s.nameResult(state, name, nil), fmt.Sprintf("update records of %s.eth", name.label), nil
}

// buildTransferTx hands a name's ERC721 registration to a sibling child
// wallet; the receiver reclaims the registry ownership as its own follow-up
// action.
func (s *Scenario) buildTransferTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, walletIdx uint64, rng *rand.Rand, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	walletCount := s.walletPool.GetWalletCount()
	targetIdx := (walletIdx + 1 + uint64(rng.Intn(int(walletCount-1)))) % walletCount
	target := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(targetIdx))
	if target == nil || target.GetAddress() == wallet.GetAddress() {
		return nil, nil, "", fmt.Errorf("no transfer target wallet available")
	}
	targetAddr := target.GetAddress()

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Base.TransferFrom(transactOpts, wallet.GetAddress(), targetAddr, new(big.Int).SetBytes(name.tokenID[:]))
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build transfer tx: %w", err)
	}

	name.inFlight = true

	onResult := func(success bool) {
		state.mu.Lock()
		name.inFlight = false
		if success {
			for i, n := range state.names {
				if n == name {
					state.names = append(state.names[:i], state.names[i+1:]...)
					break
				}
			}
		}
		state.mu.Unlock()

		if success {
			// hand the name to the receiver's state; it reclaims registry
			// ownership on its next action
			targetState := s.getWalletState(targetAddr)
			targetState.mu.Lock()
			name.needReclaim = true
			targetState.names = append(targetState.names, name)
			targetState.mu.Unlock()
		}
	}

	return tx, onResult, fmt.Sprintf("transfer %s.eth to wallet %d", name.label, targetIdx+1), nil
}

// buildAbandonTx gives up a name by zeroing its registry owner (the closest
// thing ENS has to unregistration); the wallet re-registers a fresh name later.
func (s *Scenario) buildAbandonTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Registry.SetOwner(transactOpts, name.node, common.Address{})
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build abandon tx: %w", err)
	}

	name.inFlight = true

	onResult := func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		name.inFlight = false
		if success {
			for i, n := range state.names {
				if n == name {
					state.names = append(state.names[:i], state.names[i+1:]...)
					break
				}
			}
		}
	}

	return tx, onResult, fmt.Sprintf("abandon %s.eth", name.label), nil
}

// buildReverseUpdateTx re-points the wallet's default (chain-agnostic) reverse
// record to one of its names. The addr.reverse record is left to the wallet
// naming service.
func (s *Scenario) buildReverseUpdateTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.DefaultReverse.SetName(transactOpts, name.label+".eth")
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build reverse update tx: %w", err)
	}

	return tx, s.nameResult(state, name, nil), fmt.Sprintf("set default reverse to %s.eth", name.label), nil
}

// buildWrapperApprovalTx performs the one-time base registrar approval the
// NameWrapper needs before names can be wrapped.
func (s *Scenario) buildWrapperApprovalTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Base.SetApprovalForAll(transactOpts, s.deployment.WrapperAddr, true)
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build wrapper approval tx: %w", err)
	}

	state.approveInFlight = true

	onResult := func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		state.approveInFlight = false
		if success {
			state.wrapperApproved = true
		}
	}

	return tx, onResult, "approve name wrapper", nil
}

// buildWrapTx wraps a name into the NameWrapper (ERC721 -> ERC1155).
func (s *Scenario) buildWrapTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Wrapper.WrapETH2LD(transactOpts, name.label, wallet.GetAddress(), 0, s.deployment.ResolverAddr)
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build wrap tx: %w", err)
	}

	wrapped := true
	return tx, s.nameResult(state, name, &wrapped), fmt.Sprintf("wrap %s.eth", name.label), nil
}

// buildUnwrapTx unwraps a wrapped name back to a plain ERC721 registration.
func (s *Scenario) buildUnwrapTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, name *nameState, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.Wrapper.UnwrapETH2LD(transactOpts, name.tokenID, wallet.GetAddress(), wallet.GetAddress())
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build unwrap tx: %w", err)
	}

	wrapped := false
	return tx, s.nameResult(state, name, &wrapped), fmt.Sprintf("unwrap %s.eth", name.label), nil
}

// buildChurnTx registers a short-lived name through the permissionless
// SpamRegistrarController (no commit-reveal, no minimum duration), so name
// expiry is observable within a run. Churn names are fire-and-forget.
func (s *Scenario) buildChurnTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, state *walletState, walletIdx uint64, rng *rand.Rand, feeCap, tipCap *big.Int) (*types.Transaction, func(success bool), string, error) {
	label := s.newChurnLabel(state, walletIdx)
	duration := churnMinDuration + rng.Int63n(churnMaxDuration-churnMinDuration)

	tx, err := wallet.BuildBoundTxWithEstimate(ctx, client, s.walletPool.GetTxPool(), &txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Value:     uint256.NewInt(0),
	}, func(transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return s.deployment.SpamController.Register(transactOpts, label, wallet.GetAddress(), big.NewInt(duration))
	})
	if err != nil {
		return nil, nil, "", fmt.Errorf("could not build churn register tx: %w", err)
	}

	return tx, noopResult, fmt.Sprintf("churn register %s.eth (%ds)", label, duration), nil
}

// nameResult returns an onResult callback that releases the name's in-flight
// flag and optionally updates its wrapped state on success.
func (s *Scenario) nameResult(state *walletState, name *nameState, wrappedOnSuccess *bool) func(success bool) {
	name.inFlight = true

	return func(success bool) {
		state.mu.Lock()
		defer state.mu.Unlock()

		name.inFlight = false
		if success && wrappedOnSuccess != nil {
			name.wrapped = *wrappedOnSuccess
		}
	}
}

// rentPriceWithBuffer quotes the controller rent for label/duration and adds a
// small buffer (the controller refunds overpayment).
func (s *Scenario) rentPriceWithBuffer(ctx context.Context, label string, duration *big.Int) (*uint256.Int, error) {
	price, err := s.deployment.Controller.RentPrice(&bind.CallOpts{Context: ctx}, label, duration)
	if err != nil {
		return nil, fmt.Errorf("could not quote rent price for %s: %w", label, err)
	}

	total := new(big.Int).Add(price.Base, price.Premium)
	total.Mul(total, big.NewInt(100+rentPriceBufferPercent))
	total.Div(total, big.NewInt(100))

	value, overflow := uint256.FromBig(total)
	if overflow {
		return nil, fmt.Errorf("rent price overflow for %s", label)
	}

	return value, nil
}
