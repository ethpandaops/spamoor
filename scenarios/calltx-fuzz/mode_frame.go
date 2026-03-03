package calltxfuzz

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/holiman/uint256"

	evmfuzz "github.com/ethpandaops/spamoor/scenarios/evm-fuzz"
	"github.com/ethpandaops/spamoor/spamoor"
)

// FrameMode constants for EIP-8141 frame transaction modes.
const (
	FrameModeDefault = 0 // Caller is ENTRY_POINT
	FrameModeVerify  = 1 // STATICCALL semantics, must call APPROVE
	FrameModeSender  = 2 // Caller is tx sender
)

// Frame represents a single execution frame in a Type 6 transaction.
type Frame struct {
	Mode     uint8
	Target   *common.Address // nil for contract creation
	GasLimit uint64
	Data     []byte
}

// FrameTxData represents the complete Type 6 transaction data (pre-signing).
// This is a stub — no signing or submission is implemented.
type FrameTxData struct {
	ChainID              uint64
	Nonce                uint64
	Sender               common.Address
	Frames               []Frame
	MaxPriorityFeePerGas *uint256.Int
	MaxFeePerGas         *uint256.Int
	MaxFeePerBlobGas     *uint256.Int
	BlobVersionedHashes  []common.Hash
}

// generateFrameTx generates frame transaction data but does NOT sign or submit it.
// This is a stub for future EIP-8141 support.
func (s *Scenario) generateFrameTx(txIdx uint64) *FrameTxData {
	effectiveTxID := txIdx + s.options.TxIdOffset
	rng := evmfuzz.NewDeterministicRNGWithSeed(effectiveTxID, s.seed)
	calldataGen := NewCalldataGenerator(rng, s.options.CalldataMaxSize)

	maxFrames := s.options.MaxFrames
	if maxFrames == 0 {
		maxFrames = 10
	}

	frameCount := rng.Intn(int(maxFrames)) + 1
	totalGas := s.options.GasLimit
	frames := make([]Frame, 0, frameCount)

	// Determine frame ordering pattern
	ordering := rng.Intn(4)

	for i := 0; i < frameCount; i++ {
		// Allocate gas (random distribution across frames)
		var frameGas uint64
		if i == frameCount-1 {
			frameGas = totalGas // Give remaining gas to last frame
		} else {
			frameGas = uint64(rng.Intn(int(totalGas/(uint64(frameCount-i))))) + 1
			totalGas -= frameGas
		}

		var mode uint8
		switch ordering {
		case 0: // VERIFY first, then DEFAULT/SENDER
			if i == 0 {
				mode = FrameModeVerify
			} else if rng.Float64() < 0.5 {
				mode = FrameModeSender
			} else {
				mode = FrameModeDefault
			}
		case 1: // VERIFY last
			if i == frameCount-1 {
				mode = FrameModeVerify
			} else if rng.Float64() < 0.5 {
				mode = FrameModeSender
			} else {
				mode = FrameModeDefault
			}
		case 2: // No VERIFY (invalid — should fail)
			if rng.Float64() < 0.5 {
				mode = FrameModeSender
			} else {
				mode = FrameModeDefault
			}
		default: // Multiple VERIFY (second should fail)
			if i < 2 {
				mode = FrameModeVerify
			} else if rng.Float64() < 0.5 {
				mode = FrameModeSender
			} else {
				mode = FrameModeDefault
			}
		}

		var target *common.Address
		if s.contractPool.Size() > 0 {
			addr := s.contractPool.GetRandomContract(rng)
			target = &addr
		}

		frames = append(frames, Frame{
			Mode:     mode,
			Target:   target,
			GasLimit: frameGas,
			Data:     calldataGen.Generate(),
		})
	}

	wallet := s.walletPool.GetWallet(spamoor.SelectWalletByIndex, int(txIdx))
	var sender common.Address
	if wallet != nil {
		sender = wallet.GetAddress()
	}

	ftx := &FrameTxData{
		ChainID:              s.walletPool.GetChainId().Uint64(),
		Nonce:                txIdx,
		Sender:               sender,
		Frames:               frames,
		MaxPriorityFeePerGas: uint256.NewInt(2000000000),  // 2 gwei
		MaxFeePerGas:         uint256.NewInt(20000000000), // 20 gwei
	}

	s.logger.WithField("txIdx", txIdx).Debugf(
		"generated frame tx: %d frames, ordering=%d",
		len(frames), ordering,
	)

	return ftx
}
