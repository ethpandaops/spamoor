package calltxfuzz

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/scenario"
	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txbuilder"
)

// sendCallTx sends a Type 2 DynFeeTx that calls a deployed fuzz contract.
func (s *Scenario) sendCallTx(
	ctx context.Context,
	txIdx uint64,
) (scenario.ReceiptChan, *types.Transaction, *spamoor.Client, *spamoor.Wallet, error) {
	if s.contractPool.Size() == 0 {
		return nil, nil, nil, nil, fmt.Errorf("contract pool is empty")
	}

	effectiveTxID := txIdx + s.options.TxIdOffset
	rng := evmfuzz.NewDeterministicRNGWithSeed(effectiveTxID, s.seed)
	calldataGen := NewCalldataGenerator(rng, s.options.CalldataMaxSize)

	// Pick target contract from pool
	contractAddr := s.contractPool.GetRandomContract(rng)
	calldata := calldataGen.Generate()

	// Select wallet and client
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

	// 75% chance: include a small random ETH value so called contracts
	// have balance for internal value transfers (CALL with value > 0).
	// 25%: zero value, which keeps "insufficient balance" errors possible.
	value := uint256.NewInt(0)
	if rng.Float64() < 0.75 {
		value = uint256.NewInt(uint64(0xa000 + rng.Intn(0x6000)))
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       s.options.GasLimit,
		To:        &contractAddr,
		Value:     value,
		Data:      calldata,
	})
	if err != nil {
		return nil, nil, client, wallet, err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
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
