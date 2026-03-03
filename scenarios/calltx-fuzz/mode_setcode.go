package calltxfuzz

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/sirupsen/logrus"

	"github.com/ethpandaops/spamoor/scenario"
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// sendSetCodeTx sends a Type 4 SetCodeTx with fuzzed authorization lists
// and calldata targeting delegated fuzz contracts.
func (s *Scenario) sendSetCodeTx(
	ctx context.Context,
	txIdx uint64,
) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	if s.contractPool.Size() == 0 {
		return nil, nil, nil, nil, fmt.Errorf("contract pool is empty")
	}

	effectiveTxID := txIdx + s.options.TxIdOffset
	rng := evmfuzz.NewDeterministicRNGWithSeed(effectiveTxID, s.seed)
	calldataGen := NewCalldataGenerator(rng, s.options.CalldataMaxSize)

	// Build authorization list
	authList, delegators := s.buildFuzzedAuthList(rng, txIdx)

	// EIP-7702 requires at least one authorization; fall back to call tx if empty
	if len(authList) == 0 {
		return s.sendCallTx(ctx, txIdx)
	}

	// Pick a delegated EOA as the To address
	var toAddr common.Address
	if len(delegators) > 0 {
		toAddr = delegators[rng.Intn(len(delegators))].GetAddress()
	} else {
		// Fallback to a pool contract
		toAddr = s.contractPool.GetRandomContract(rng)
	}

	calldata := calldataGen.Generate()

	// Select sender wallet and client
	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByPendingTxCount, int(txIdx))
	if wallet == nil {
		return nil, nil, nil, nil, fmt.Errorf("no wallet available")
	}

	client := s.walletPool.GetClient(
		spamoor.WithClientSelectionMode(spamoor.SelectClientByIndex, int(txIdx)),
		spamoor.WithClientGroup(s.options.ClientGroup),
	)
	if client == nil {
		return nil, nil, nil, wallet, fmt.Errorf("no client available")
	}

	if err := wallet.ResetNoncesIfNeeded(ctx, client); err != nil {
		return nil, nil, client, wallet, err
	}

	baseFeeWei, tipFeeWei := spamoor.ResolveFees(
		s.options.BaseFee, s.options.TipFee,
		s.options.BaseFeeWei, s.options.TipFeeWei,
	)
	feeCap, tipCap, err := s.walletPool.GetTxPool().GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	// 75% chance: include a small random ETH value so delegated contracts
	// have balance for internal value transfers (CALL with value > 0).
	// 25%: zero value, which keeps "insufficient balance" errors possible.
	value := uint256.NewInt(0)
	if rng.Float64() < 0.75 {
		value = uint256.NewInt(uint64(0xa000 + rng.Intn(0x6000)))
	}

	txData, err := txbuilder.SetCodeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &toAddr,
		Value:     value,
		Data:      calldata,
		AuthList:  authList,
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := wallet.BuildSetCodeTx(txData)
	if err != nil {
		return nil, nil, client, wallet, err
	}

	receiptChan := make(scenario.ReceiptChan, 1)

	err = s.walletPool.GetTxPool().SendTransaction(ctx, wallet, tx, &spamoor.SendTransactionOptions{
		Client:      client,
		ClientGroup: s.options.ClientGroup,
		Rebroadcast: s.options.Rebroadcast > 0,
		OnComplete: func(tx *types.Transaction, receipt *types.Receipt, err error) {
			receiptChan <- receipt
		},
	})
	if err != nil {
		wallet.MarkSkippedNonce(tx.Nonce())
		return nil, nil, client, wallet, err
	}

	return receiptChan, tx, client, wallet, nil
}

// buildFuzzedAuthList constructs a fuzzed authorization list for SetCodeTx.
// Returns the authorization list and the delegator wallets used.
func (s *Scenario) buildFuzzedAuthList(
	rng *evmfuzz.DeterministicRNG,
	txIdx uint64,
) ([]types.SetCodeAuthorization, []*spamoor.Wallet) {
	minAuth := s.options.MinAuthorizations
	maxAuth := s.options.MaxAuthorizations
	if maxAuth == 0 {
		return nil, nil
	}
	if minAuth > maxAuth {
		minAuth = maxAuth
	}

	n := big.NewInt(int64(maxAuth - minAuth + 1))
	authCount := int(minAuth) + rng.Intn(int(n.Int64()))

	authorizations := make([]types.SetCodeAuthorization, 0, authCount)
	delegators := make([]*spamoor.Wallet, 0, authCount)
	chainID := s.walletPool.GetChainId().Uint64()

	for i := 0; i < authCount; i++ {
		delegatorIdx := (txIdx * maxAuth) + uint64(i)
		if s.options.MaxDelegators > 0 {
			delegatorIdx = delegatorIdx % s.options.MaxDelegators
		}

		delegator := s.getOrCreateDelegator(delegatorIdx)
		if delegator == nil {
			continue
		}
		delegators = append(delegators, delegator)

		// Decide if this authorization should be deliberately invalid
		if rng.Float64() < s.options.InvalidAuthRatio {
			auth := s.buildInvalidAuth(rng, delegator, chainID)
			authorizations = append(authorizations, auth)
			continue
		}

		// Valid authorization: delegate to a pool contract
		codeAddr := s.contractPool.GetRandomContract(rng)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}

		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			s.logger.WithFields(logrus.Fields{
				"delegator": delegator.GetAddress().Hex(),
			}).Warnf("failed to sign auth: %v", err)
			continue
		}
		authorizations = append(authorizations, signedAuth)
	}

	return authorizations, delegators
}

// buildInvalidAuth constructs a deliberately invalid authorization for fuzzing.
func (s *Scenario) buildInvalidAuth(
	rng *evmfuzz.DeterministicRNG,
	delegator *spamoor.Wallet,
	chainID uint64,
) types.SetCodeAuthorization {
	invalidType := rng.Intn(7)

	switch invalidType {
	case 0: // Wrong chain_id
		codeAddr := s.contractPool.GetRandomContract(rng)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(uint64(rng.Intn(1000000) + 1)),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth

	case 1: // Chain_id = 0 (wildcard)
		codeAddr := s.contractPool.GetRandomContract(rng)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(0),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth

	case 2: // Invalid signature (corrupted R/S)
		codeAddr := s.contractPool.GetRandomContract(rng)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		// Corrupt the signature
		signedAuth.R.SetUint64(rng.Uint64())
		return signedAuth

	case 3: // Nonce mismatch
		codeAddr := s.contractPool.GetRandomContract(rng)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: codeAddr,
			Nonce:   delegator.GetNextNonce() + uint64(rng.Intn(100)+1),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth

	case 4: // Delegation to address(0) (clear delegation)
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: common.Address{},
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth

	case 5: // Delegation to precompile address
		precompileAddr := common.BigToAddress(big.NewInt(int64(rng.Intn(9) + 1)))
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: precompileAddr,
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth

	default: // Delegation to another delegator EOA (chain test)
		var targetAddr common.Address
		otherIdx := uint64(rng.Intn(int(s.options.MaxDelegators) + 1))
		otherDelegator := s.getOrCreateDelegator(otherIdx)
		if otherDelegator != nil {
			targetAddr = otherDelegator.GetAddress()
		} else {
			targetAddr = common.BigToAddress(big.NewInt(int64(rng.Uint64())))
		}
		auth := types.SetCodeAuthorization{
			ChainID: *uint256.NewInt(chainID),
			Address: targetAddr,
			Nonce:   delegator.GetNextNonce(),
		}
		signedAuth, err := types.SignSetCode(delegator.GetPrivateKey(), auth)
		if err != nil {
			return auth
		}
		return signedAuth
	}
}

// getOrCreateDelegator returns or creates a delegator wallet at the given index.
func (s *Scenario) getOrCreateDelegator(idx uint64) *spamoor.Wallet {
	s.delegatorMu.Lock()
	defer s.delegatorMu.Unlock()

	if s.options.MaxDelegators > 0 && len(s.delegators) > int(idx) && s.delegators[idx] != nil {
		return s.delegators[idx]
	}

	delegator, err := s.prepareDelegator(idx)
	if err != nil {
		s.logger.Errorf("could not prepare delegator %d: %v", idx, err)
		return nil
	}

	if s.options.MaxDelegators > 0 {
		// Grow slice if needed
		for len(s.delegators) <= int(idx) {
			s.delegators = append(s.delegators, nil)
		}
		s.delegators[idx] = delegator
	}

	return delegator
}

// prepareDelegator derives a deterministic delegator wallet.
func (s *Scenario) prepareDelegator(idx uint64) (*spamoor.Wallet, error) {
	idxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(idxBytes, idx)

	if s.options.MaxDelegators > 0 {
		idxBytes = append(idxBytes, s.delegatorSeed...)
	}

	rootAddr := s.walletPool.GetRootWallet().GetWallet().GetAddress()
	childKey := sha256.Sum256(append(common.FromHex(rootAddr.Hex()), idxBytes...))

	privateKey, address, err := spamoor.LoadPrivateKey(fmt.Sprintf("%x", childKey))
	if err != nil {
		return nil, err
	}

	return spamoor.NewWallet(privateKey, address), nil
}
