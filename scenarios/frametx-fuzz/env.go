package frametxfuzz

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// environment is everything the generator needs from the chain.
//
// It is assembled once, before any transaction is generated, so that generation itself
// stays a pure function of the recipe. Anything discovered here -- which extensions the
// chain runs, where the probe contract lives, which wallets carry the delegation, what
// the current slot is -- is applied to a recipe rather than drawn from it.
type environment struct {
	logger      logrus.FieldLogger
	walletPool  *spamoor.WalletPool
	clientGroup string

	// extensions is the envelope shape the chain expects.
	extensions txtypes.FrameExtensions

	// probe is the deployed probe contract, or nil when probe axes are disabled.
	probe *ProbeDeployment

	// plainFrom and contractFrom bound the two halves of the wallet pool. A wallet
	// carrying the probe delegation has code, so its validation frame runs the probe
	// contract instead of the protocol's default code; mixing the two would silently
	// change what a default-code recipe tests.
	plainCount    int
	contractCount int

	// nonces tracks the EIP-8250 keyed nonce sequences this run has consumed.
	nonces *nonceLedger

	// roots tracks the EIP-8272 roots this run has committed.
	roots *rootRing

	// burner sends the transactions that are not meant to land.
	burner *spamoor.Wallet

	// allowPostTx and allowBlobs record what the chain accepted during the startup
	// probes, so a recipe asking for something the chain refuses is downgraded rather
	// than reported as a finding on every transaction.
	allowPostTx bool
	allowBlobs  bool
}

// plainWallet returns a wallet with no code, for recipes validated by default code.
func (e *environment) plainWallet(index int) *spamoor.Wallet {
	if e.plainCount == 0 {
		return nil
	}

	return e.walletPool.GetWallet(spamoor.SelectWalletByIndex, index%e.plainCount)
}

// contractWallet returns a wallet carrying the probe delegation.
func (e *environment) contractWallet(index int) *spamoor.Wallet {
	if e.contractCount == 0 {
		return nil
	}

	return e.walletPool.GetWallet(spamoor.SelectWalletByIndex, e.plainCount+index%e.contractCount)
}

// senderFor picks the wallet a recipe calls for.
//
// A recipe carrying a deliberate violation is sent from the burner: it will never land,
// and drawing a nonce from a generating wallet for a transaction that never confirms
// would stall everything that wallet sends afterwards.
func (e *environment) senderFor(recipe *Recipe, index int) *spamoor.Wallet {
	if recipe.Invalid != "" && e.burner != nil {
		return e.burner
	}

	if recipe.Sender == SenderContract && e.contractCount > 0 {
		return e.contractWallet(index)
	}

	return e.plainWallet(index)
}

// targetWallet returns a funded wallet to address, which is never the sender.
func (e *environment) targetWallet(index int) common.Address {
	wallet := e.plainWallet(index + 1)
	if wallet == nil {
		return common.Address{}
	}

	return wallet.GetAddress()
}

// setupEnvironment probes the chain and deploys what the enabled axes need.
func (s *Scenario) setupEnvironment(ctx context.Context) (*environment, error) {
	txpool := s.walletPool.GetTxPool()
	if txpool == nil {
		return nil, fmt.Errorf("the scenario has no transaction pool")
	}

	if !txpool.IsAmsterdam() {
		return nil, fmt.Errorf("frame transactions need the Amsterdam (EIP-8037) gas model, but --pre-amsterdam-fee-model is set")
	}

	support, err := txpool.GetFrameSupportWithInit(ctx)
	if err != nil {
		return nil, err
	}

	if !support.Active {
		return nil, fmt.Errorf("no account at the EIP-8141 expiry verifier predeploy %s: this chain does not implement frame transactions",
			txtypes.ExpiryVerifier)
	}

	env := &environment{
		logger:      s.logger,
		walletPool:  s.walletPool,
		clientGroup: s.options.ClientGroup,
		extensions:  support.Extensions,
	}

	if s.pinnedExtensions != nil {
		env.extensions = *s.pinnedExtensions
		s.logger.Infof("using pinned frame transaction envelope: %s", env.extensions)
	} else {
		s.logger.Infof("detected frame transaction envelope: %s", env.extensions)
	}

	total := int(s.walletPool.GetWalletCount())
	env.plainCount = total
	env.burner = s.walletPool.GetWellKnownWallet(BurnerWalletName)

	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	feeCap, tipCap, err := s.fees(client)
	if err != nil {
		return nil, err
	}

	if s.axes.enabled(axisProbe) {
		if err := s.setupProbe(ctx, env, total, feeCap, tipCap); err != nil {
			return nil, err
		}
	}

	if env.extensions.Has(txtypes.FrameExtKeyedNonces) && s.axes.enabled(axisNonces) {
		env.nonces = newNonceLedger()
	}

	if env.extensions.Has(txtypes.FrameExtRecentRoots) && s.axes.enabled(axisRoots) {
		env.roots, err = newRootRing(ctx, s.logger, client)
		if err != nil {
			s.logger.Warnf("recent root axis disabled: %v", err)

			env.roots = nil
		}
	}

	return env, nil
}

// setupProbe deploys the probe contract and splits the wallet pool.
func (s *Scenario) setupProbe(ctx context.Context, env *environment, total int, feeCap, tipCap *big.Int) error {
	deployment, err := SetupProbe(ctx, s.logger, s.walletPool, s.options.ClientGroup, feeCap, tipCap)
	if err != nil {
		return err
	}

	env.probe = deployment

	// The paymaster approves payment from its own code, so it needs the probe contract
	// delegated to it exactly as a contract sender does.
	paymaster := s.walletPool.GetWellKnownWallet(ProbeDeployerWalletName)
	if err := DelegateProbe(ctx, s.logger, s.walletPool, s.options.ClientGroup, deployment.Address,
		[]*spamoor.Wallet{paymaster}, feeCap, tipCap); err != nil {
		return err
	}

	if total < 4 {
		s.logger.Warnf("wallet pool of %d is too small to reserve delegated senders; contract-sender recipes will use default code", total)

		return nil
	}

	env.contractCount = max(total/4, 1)
	env.plainCount = total - env.contractCount

	wallets := make([]*spamoor.Wallet, 0, env.contractCount)
	for i := 0; i < env.contractCount; i++ {
		wallets = append(wallets, s.walletPool.GetWallet(spamoor.SelectWalletByIndex, env.plainCount+i))
	}

	return DelegateProbe(ctx, s.logger, s.walletPool, s.options.ClientGroup, deployment.Address, wallets, feeCap, tipCap)
}

// fees resolves the configured fee caps.
func (s *Scenario) fees(client *spamoor.Client) (*big.Int, *big.Int, error) {
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)

	return s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
}

// probeCapabilities runs the A/B probes for the features no predeploy announces.
//
// EIP-7906 installs nothing, and a frame transaction carrying blobs is refused outright
// by at least one client build. Both are settled by sending two transactions that differ
// only in the feature under test: if the plain one lands and the other does not, the
// feature is unsupported. The signal is the accept/reject of a controlled pair, never
// the text of an error.
func (s *Scenario) probeCapabilities(ctx context.Context, env *environment) {
	if s.options.PostTx == "off" || !s.axes.enabled(axisPostTx) {
		env.allowPostTx = false
	} else if s.options.PostTx == "on" {
		env.allowPostTx = true
	} else {
		env.allowPostTx = s.abProbe(ctx, env, "POST_TX", func(frames []*txtypes.Frame) []*txtypes.Frame {
			return append(frames, txtypes.PostTxFrame(txtypes.ExpiryVerifier, expiryData(s.deadline()),
				txtypes.FrameLimits{Execution: s.options.VerifyGas}))
		})
	}

	if s.options.PostTx == "auto" && s.axes.enabled(axisPostTx) {
		s.logger.Infof("EIP-7906 POST_TX frames are %s on this chain", supportWord(env.allowPostTx))
	}
}

// supportWord renders a capability probe's result.
func supportWord(supported bool) string {
	if supported {
		return "supported"
	}

	return "not supported"
}

// abProbe sends a minimal transaction with and without a feature and reports whether the
// chain accepted both.
func (s *Scenario) abProbe(ctx context.Context, env *environment, name string, addFeature func([]*txtypes.Frame) []*txtypes.Frame) bool {
	client := s.walletPool.GetClient(spamoor.WithClientGroup(s.options.ClientGroup))
	if client == nil {
		return false
	}

	wallet := env.plainWallet(0)
	if wallet == nil {
		return false
	}

	feeCap, tipCap, err := s.fees(client)
	if err != nil {
		return false
	}

	baseline := func() []*txtypes.Frame {
		return []*txtypes.Frame{
			txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: s.options.VerifyGas}),
			txtypes.UserOpFrame(nil, nil, nil, txtypes.FrameLimits{Execution: s.options.UserOpGas}),
		}
	}

	if err := s.sendProbeTx(ctx, client, wallet, env, baseline(), feeCap, tipCap); err != nil {
		s.logger.Warnf("could not establish a baseline for the %s probe, assuming unsupported: %v", name, err)

		return false
	}

	if err := s.sendProbeTx(ctx, client, wallet, env, addFeature(baseline()), feeCap, tipCap); err != nil {
		s.logger.Debugf("%s probe transaction was refused: %v", name, err)

		return false
	}

	return true
}

// sendProbeTx submits one capability probe transaction and waits for its receipt.
func (s *Scenario) sendProbeTx(ctx context.Context, client *spamoor.Client, wallet *spamoor.Wallet, env *environment, frames []*txtypes.Frame, feeCap, tipCap *big.Int) error {
	frameTx := txtypes.NewFrameTxWithExtensions(env.extensions, nil, wallet.GetAddress(), 0,
		txtypes.FrameFees{
			GasTipCap: uint256.MustFromBig(tipCap),
			GasFeeCap: uint256.MustFromBig(feeCap),
		},
		frames,
		[]*txtypes.FrameSignature{txtypes.SenderSignature()},
	)

	tx, err := wallet.BuildFrameTx(frameTx)
	if err != nil {
		return err
	}

	if _, err := s.walletPool.GetTxPool().SendAndAwaitTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
	}); err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())

		return err
	}

	return nil
}
