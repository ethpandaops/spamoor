package frametxfuzz

import (
	"context"
	"fmt"
	"math/big"

	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

// environment is everything the generator needs from the chain. It is assembled once,
// before any transaction is generated, so generation stays a function of the recipe.
type environment struct {
	logger      logrus.FieldLogger
	walletPool  *spamoor.WalletPool
	clientGroup string

	// extensions is the envelope shape the chain expects.
	extensions txtypes.FrameExtensions

	// probe is the deployed probe contract, or nil when probe axes are disabled.
	probe *ProbeDeployment

	// The wallet pool is split in two. A wallet carrying the probe delegation has code,
	// so its validation frame runs that code instead of the protocol's default code.
	plainCount    int
	contractCount int

	// nonces tracks the EIP-8250 keyed nonce sequences this run has consumed.
	nonces *nonceLedger

	// roots tracks the EIP-8272 roots this run has committed.
	roots *rootRing

	// burner sends the transactions that are not meant to land.
	burner *spamoor.Wallet

	// factory is the CREATE2 factory a deploy frame calls, with a salt followed by init
	// code. The address is computable before sending, so a later frame can name it.
	factory common.Address

	// contracts are the generated contracts earlier transactions deployed, so a frame
	// can call code that some other transaction's frame created.
	contractsMutex sync.Mutex
	contracts      []common.Address

	// allowPostTx and allowBlobs record what the startup probes found, so a recipe
	// asking for something the chain refuses is downgraded rather than sent.
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

// senderFor picks the wallet a recipe calls for. A recipe carrying a deliberate
// violation goes to the burner, since it never lands and would stall a real wallet.
func (e *environment) senderFor(recipe *Recipe, index int) *spamoor.Wallet {
	if recipe.Invalid != "" && e.burner != nil {
		return e.burner
	}

	if recipe.Sender == SenderContract && e.contractCount > 0 {
		return e.contractWallet(index)
	}

	return e.plainWallet(index)
}

// rememberContracts records generated contracts a landed transaction deployed, keeping
// the most recent ones.
func (e *environment) rememberContracts(addresses []common.Address) {
	const keep = 256

	if len(addresses) == 0 {
		return
	}

	e.contractsMutex.Lock()
	defer e.contractsMutex.Unlock()

	e.contracts = append(e.contracts, addresses...)
	if len(e.contracts) > keep {
		e.contracts = e.contracts[len(e.contracts)-keep:]
	}
}

// pickContract returns a previously deployed generated contract, if there is one.
func (e *environment) pickContract(index int) (common.Address, bool) {
	e.contractsMutex.Lock()
	defer e.contractsMutex.Unlock()

	if len(e.contracts) == 0 {
		return common.Address{}, false
	}

	return e.contracts[index%len(e.contracts)], true
}

// paymasterFor picks the account that sponsors a transaction, spread across every wallet
// carrying the probe delegation.
//
// The public mempool caps how many pending transactions one non-canonical paymaster may
// sponsor, so a single paymaster refuses almost everything behind the first few. The
// sender is skipped: a transaction sponsoring itself is the self-relayed shape.
func (e *environment) paymasterFor(index int, sender common.Address) common.Address {
	if e.probe == nil {
		return common.Address{}
	}

	for offset := 0; offset < e.contractCount; offset++ {
		wallet := e.contractWallet(index + offset)
		if wallet != nil && wallet.GetAddress() != sender {
			return wallet.GetAddress()
		}
	}

	return e.probe.Paymaster
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

	if s.axes.enabled(axisCode) {
		env.factory, err = s.walletPool.GetDeploymentFactory().GetFactoryAddress(ctx)
		if err != nil {
			s.logger.Warnf("generated contract axis disabled: %v", err)
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

	// The paymaster approves payment from its own code, so it needs the delegation too.
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

// probeCapabilities settles the features no predeploy announces by sending two
// transactions that differ only in the feature: if the plain one lands and the other does
// not, the feature is unsupported. The signal is accept/reject, never error text.
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
