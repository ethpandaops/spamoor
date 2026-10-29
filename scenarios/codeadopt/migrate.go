package codeadopt

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/holiman/uint256"

	"github.com/ethpandaops/spamoor/spamoor"
	"github.com/ethpandaops/spamoor/txtypes"
)

const (
	approveVerifyGas  = 20_000
	frameTransferGas  = 30_000
	frameReceiptAwait = 3 * time.Minute
)

func ecrecoverInput(key *ecdsa.PrivateKey) ([]byte, error) {
	digest := make([]byte, 32)
	copy(digest, randomAddress().Bytes())
	copy(digest[20:], randomAddress().Bytes())

	sig, err := crypto.Sign(digest, key)
	if err != nil {
		return nil, err
	}

	input := make([]byte, 0, 128)
	input = append(input, digest...)
	input = append(input, common.LeftPadBytes([]byte{sig[64] + 27}, 32)...)
	input = append(input, sig[:64]...)

	return input, nil
}

func (s *Scenario) fees(client *spamoor.Client) (*uint256.Int, *uint256.Int, error) {
	baseFeeWei, tipFeeWei := spamoor.ResolveFees(s.options.BaseFee, s.options.TipFee, s.options.BaseFeeWei, s.options.TipFeeWei)

	feeCap, tipCap, err := s.walletPool.GetSuggestedFees(client, baseFeeWei, tipFeeWei)
	if err != nil {
		return nil, nil, err
	}

	return uint256.MustFromBig(feeCap), uint256.MustFromBig(tipCap), nil
}

func (s *Scenario) sendECDSA(ctx context.Context, client *spamoor.Client, key *ecdsa.PrivateKey, nonce uint64, to common.Address) error {
	feeCap, tipCap, err := s.fees(client)
	if err != nil {
		return fmt.Errorf("getting fees: %w", err)
	}

	tx, err := types.SignNewTx(key, types.LatestSignerForChainID(s.walletPool.GetChainId()), &types.DynamicFeeTx{
		ChainID:   s.walletPool.GetChainId(),
		Nonce:     nonce,
		GasTipCap: tipCap.ToBig(),
		GasFeeCap: feeCap.ToBig(),
		Gas:       50_000,
		To:        &to,
		Value:     common.Big1,
	})
	if err != nil {
		return fmt.Errorf("signing: %w", err)
	}

	raw, err := tx.MarshalBinary()
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}

	return client.SendRawTransaction(ctx, raw)
}

func (s *Scenario) sendFrameFrom(ctx context.Context, client *spamoor.Client, sender common.Address, nonce uint64, to common.Address) (*txtypes.Receipt, error) {
	feeCap, tipCap, err := s.fees(client)
	if err != nil {
		return nil, fmt.Errorf("getting fees: %w", err)
	}

	verify := txtypes.SelfVerifyFrame(txtypes.FrameLimits{Execution: approveVerifyGas})
	verify.Data = common.LeftPadBytes([]byte{verify.ApprovalScope()}, 32)

	transfer := txtypes.UserOpFrame(&to, uint256.NewInt(1), nil, txtypes.FrameLimits{Execution: frameTransferGas})

	frameTx := txtypes.NewFrameTxWithExtensions(s.support.Extensions, uint256.MustFromBig(s.walletPool.GetChainId()), sender, 0,
		txtypes.FrameFees{GasTipCap: tipCap, GasFeeCap: feeCap},
		[]*txtypes.Frame{verify, transfer},
		[]*txtypes.FrameSignature{},
	)

	if frameTx.HasKeyedNonces() && len(frameTx.NonceKeys) == 0 {
		frameTx.WithLegacyNonce()
	}

	if frameTx.UsesLegacyNonce() {
		frameTx.NonceSeq = nonce
	}

	if err := frameTx.DeriveMaxFee(); err != nil {
		return nil, fmt.Errorf("deriving max fee: %w", err)
	}

	tx := txtypes.NewTx(frameTx)

	raw, err := tx.MarshalNetwork()
	if err != nil {
		return nil, fmt.Errorf("encoding: %w", err)
	}

	if err := client.SendRawTransaction(ctx, raw); err != nil {
		return nil, fmt.Errorf("sending: %w", err)
	}

	waitCtx, cancel := context.WithTimeout(ctx, frameReceiptAwait)
	defer cancel()

	for {
		receipt, err := client.GetTransactionReceipt(waitCtx, tx.Hash())
		if err == nil && receipt != nil {
			return receipt, nil
		}

		select {
		case <-waitCtx.Done():
			return nil, fmt.Errorf("frame transaction %s not included: %w", tx.Hash(), waitCtx.Err())
		case <-time.After(2 * time.Second):
		}
	}
}
