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

const (
	// ProbeDeployerWalletName is the well-known wallet the probe contract is deployed
	// from, and the account that plays the paymaster.
	ProbeDeployerWalletName = "frametx-fuzz-probe"

	// ProbeDelegatorWalletName is the well-known wallet that carries delegation
	// transactions. It exists so that no wallet ever delegates itself: an EIP-7702
	// authorization signed by the transaction's own sender must use the nonce after
	// the transaction's, and mixing the two sequences in one wallet leaves a nonce gap
	// that stalls everything else it sends.
	ProbeDelegatorWalletName = "frametx-fuzz-delegator"
)

// ProbeDeployment is the result of setting the probe contract up on a chain.
type ProbeDeployment struct {
	// Address is where the probe contract lives. It is deterministic, so a rerun finds
	// the same address and skips deployment.
	Address common.Address

	// Paymaster is the account a pay frame targets. It is the deployer wallet with the
	// probe contract delegated to it, so its own code runs APPROVE(APPROVE_PAYMENT).
	Paymaster common.Address
}

// SetupProbe deploys the probe contract if it is not already present.
//
// The contract is deployed through the CREATE2 factory with a fixed salt, which makes
// the address a function of the code alone: every scenario, every run and every restart
// addresses the same contract, and a chain that already has it costs nothing.
func SetupProbe(ctx context.Context, logger logrus.FieldLogger, walletPool *spamoor.WalletPool, clientGroup string, feeCap, tipCap *big.Int) (*ProbeDeployment, error) {
	client := walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(clientGroup),
	)
	if client == nil {
		return nil, scenario.ErrNoClients
	}

	deployer := walletPool.GetWellKnownWallet(ProbeDeployerWalletName)
	if deployer == nil {
		return nil, fmt.Errorf("the %q well-known wallet is not registered", ProbeDeployerWalletName)
	}

	initCode, err := ProbeInitCode()
	if err != nil {
		return nil, err
	}

	// A fixed salt: the address should depend on the contract, not on who deployed it.
	salt := [32]byte{}
	copy(salt[:], []byte("spamoor-frame-probe-v1"))

	address, tx, err := walletPool.GetDeploymentFactory().GetContractDeployment(ctx, initCode, salt, client, deployer, feeCap, tipCap, false)
	if err != nil {
		return nil, fmt.Errorf("could not prepare the frame probe deployment: %w", err)
	}

	deployment := &ProbeDeployment{
		Address:   address,
		Paymaster: deployer.GetAddress(),
	}

	if tx == nil {
		logger.Infof("frame probe contract already deployed at %v", address)

		return deployment, nil
	}

	if _, err := walletPool.GetTxPool().SendAndAwaitTransaction(ctx, deployer, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: clientGroup,
		Rebroadcast: true,
	}); err != nil {
		return nil, fmt.Errorf("could not deploy the frame probe contract: %w", err)
	}

	logger.Infof("deployed frame probe contract at %v", address)

	return deployment, nil
}

// authorizationsPerTx bounds how many delegations ride in one transaction, keeping the
// gas well inside the per-transaction cap.
const authorizationsPerTx = 16

// delegationGasPerAuthorization is a deliberately generous budget for one authorization.
// Under EIP-8037 an authorization's cost is split across both gas dimensions, and a
// delegation installed once per run is not worth being tight about.
const delegationGasPerAuthorization = 80_000

// DelegateProbe installs an EIP-7702 delegation to the probe contract on each wallet, giving
// them code without a deploy frame.
//
// This is the only practical way to exercise a contract sender. The alternative EIP-8141
// offers is a deploy frame that CREATE2s account code at tx.sender, but a CREATE2 account
// deployment costs roughly 224,000 gas against the 100,000 execution cap that applies to
// the whole validation prefix, so it cannot propagate through the public mempool. A
// delegation is installed once by an ordinary transaction and costs the frame nothing.
//
// Each wallet signs its own authorization while a separate sender carries them, so the
// authorization's nonce is drawn from the wallet whose nonce the protocol will actually
// increment. Wallets that already carry the delegation are skipped, so a rerun is free.
func DelegateProbe(ctx context.Context, logger logrus.FieldLogger, walletPool *spamoor.WalletPool, clientGroup string, target common.Address, wallets []*spamoor.Wallet, feeCap, tipCap *big.Int) error {
	client := walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, 0),
		spamoor.WithClientGroup(clientGroup),
	)
	if client == nil {
		return scenario.ErrNoClients
	}

	sender := walletPool.GetWellKnownWallet(ProbeDelegatorWalletName)
	if sender == nil {
		return fmt.Errorf("the %q well-known wallet is not registered", ProbeDelegatorWalletName)
	}

	for _, wallet := range wallets {
		if wallet != nil && wallet.GetAddress() == sender.GetAddress() {
			return fmt.Errorf("the delegation sender %v cannot delegate itself", sender.GetAddress())
		}
	}

	pending, err := undelegatedWallets(ctx, client, target, wallets)
	if err != nil {
		return err
	}

	if len(pending) == 0 {
		return nil
	}

	logger.Infof("installing the frame probe delegation on %v wallets", len(pending))

	for start := 0; start < len(pending); start += authorizationsPerTx {
		end := min(start+authorizationsPerTx, len(pending))

		if err := delegateBatch(ctx, walletPool, clientGroup, client, sender, target, pending[start:end], feeCap, tipCap); err != nil {
			return err
		}
	}

	return nil
}

// undelegatedWallets returns the wallets that do not already delegate to target.
func undelegatedWallets(ctx context.Context, client *spamoor.Client, target common.Address, wallets []*spamoor.Wallet) ([]*spamoor.Wallet, error) {
	pending := make([]*spamoor.Wallet, 0, len(wallets))

	for _, wallet := range wallets {
		if wallet == nil {
			continue
		}

		code, err := client.GetCodeAt(ctx, wallet.GetAddress())
		if err != nil {
			return nil, fmt.Errorf("could not read code of %v: %w", wallet.GetAddress(), err)
		}

		if delegate, ok := txtypes.ParseDelegation(code); ok && delegate == target {
			continue
		}

		pending = append(pending, wallet)
	}

	return pending, nil
}

// delegateBatch installs the delegation on one chunk of wallets in a single transaction.
func delegateBatch(ctx context.Context, walletPool *spamoor.WalletPool, clientGroup string, client *spamoor.Client, sender *spamoor.Wallet, target common.Address, wallets []*spamoor.Wallet, feeCap, tipCap *big.Int) error {
	chainID := uint256.MustFromBig(walletPool.GetChainId())
	authorizations := make([]txtypes.SetCodeAuthorization, 0, len(wallets))

	for _, wallet := range wallets {
		key := wallet.GetPrivateKey()
		if key == nil {
			return fmt.Errorf("wallet %v has no private key to authorize a delegation with", wallet.GetAddress())
		}

		// The protocol increments the authority's nonce when it applies the
		// authorization, so the nonce has to come out of that wallet's own sequence.
		authorization, err := txtypes.SignAuthorization(txtypes.SetCodeAuthorization{
			ChainID: *chainID,
			Address: target,
			Nonce:   wallet.GetNextNonce(),
		}, key)
		if err != nil {
			return fmt.Errorf("could not sign the delegation for %v: %w", wallet.GetAddress(), err)
		}

		authorizations = append(authorizations, authorization)
	}

	tx, err := sender.BuildSetCodeTx(&txtypes.SetCodeTx{
		GasTipCap: uint256.MustFromBig(tipCap),
		GasFeeCap: uint256.MustFromBig(feeCap),
		Gas:       uint64(len(authorizations)) * delegationGasPerAuthorization,
		To:        sender.GetAddress(),
		Value:     uint256.NewInt(0),
		AuthList:  authorizations,
	})
	if err != nil {
		return fmt.Errorf("could not build the delegation transaction: %w", err)
	}

	if _, err := walletPool.GetTxPool().SendAndAwaitTransaction(ctx, sender, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: clientGroup,
		Rebroadcast: true,
	}); err != nil {
		return fmt.Errorf("could not install delegations: %w", err)
	}

	return nil
}
