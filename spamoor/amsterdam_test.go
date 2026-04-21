package spamoor

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func TestComputeCostPerStateByte(t *testing.T) {
	// Reference anchors from the EIP-8037 specification / go-ethereum
	// glamsterdam-devnet-0 hardcoded value.
	tests := []struct {
		name     string
		gasLimit uint64
		want     uint64
	}{
		{name: "zero gas limit", gasLimit: 0, want: 0},
		// Any non-zero gas limit below the 100M anchor gets clamped to the
		// devnet-0 floor of 1174 (geth hardcodes this, see cpsbFloor).
		{name: "tiny gas limit clamps to floor", gasLimit: 1, want: cpsbFloor},
		{name: "60M gas limit clamps to floor", gasLimit: 60_000_000, want: cpsbFloor},
		{name: "100M gas limit matches formula", gasLimit: 100_000_000, want: 1174},
		{name: "200M gas limit uses formula above floor", gasLimit: 200_000_000, want: 2198},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, computeCostPerStateByte(tt.gasLimit))
		})
	}
}

// buildTestPool constructs a minimally-initialized TxPool whose only wired
// state is what FundingGasFor / batcherGasFor read: the Amsterdam flag and
// the current gas limit (which drives cpsb). Avoids pulling in the full
// NewTxPool bootstrap (RPC, goroutines, block poller).
func buildTestPool(isAmsterdam bool, gasLimit uint64) *TxPool {
	return &TxPool{isAmsterdam: isAmsterdam, currentGasLimit: gasLimit}
}

func buildTestWalletPool(isAmsterdam bool, gasLimit, fundingOverride uint64) *WalletPool {
	return &WalletPool{
		txpool: buildTestPool(isAmsterdam, gasLimit),
		config: WalletPoolConfig{FundingGasLimit: fundingOverride},
	}
}

func TestFundingGasFor(t *testing.T) {
	tests := []struct {
		name        string
		isAmsterdam bool
		gasLimit    uint64
		override    uint64
		isEmpty     bool
		want        uint64
	}{
		{
			name: "pre-amsterdam empty target", isAmsterdam: false, gasLimit: 60_000_000,
			override: 0, isEmpty: true, want: 21_000,
		},
		{
			name: "pre-amsterdam non-empty target", isAmsterdam: false, gasLimit: 60_000_000,
			override: 0, isEmpty: false, want: 21_000,
		},
		{
			name: "amsterdam non-empty target 60M", isAmsterdam: true, gasLimit: 60_000_000,
			override: 0, isEmpty: false, want: 21_000 * 10 / 9, // 23,100
		},
		{
			// cpsb is clamped to the floor (1174) at 60M gas limit, matching
			// the value devnet-0 geth hardcodes.
			name: "amsterdam empty target 60M (cpsb clamped to floor)", isAmsterdam: true, gasLimit: 60_000_000,
			override: 0, isEmpty: true, want: (21_000 + 112*cpsbFloor) * 10 / 9,
		},
		{
			name: "amsterdam empty target 100M (cpsb=1174)", isAmsterdam: true, gasLimit: 100_000_000,
			override: 0, isEmpty: true, want: (21_000 + 112*1174) * 10 / 9, // 167,736
		},
		{
			name: "override wins pre-amsterdam", isAmsterdam: false, gasLimit: 60_000_000,
			override: 42_000, isEmpty: true, want: 42_000,
		},
		{
			name: "override wins amsterdam", isAmsterdam: true, gasLimit: 100_000_000,
			override: 42_000, isEmpty: true, want: 42_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := buildTestWalletPool(tt.isAmsterdam, tt.gasLimit, tt.override)
			require.Equal(t, tt.want, pool.FundingGasFor(tt.isEmpty))
		})
	}
}

func TestBatcherGasFor(t *testing.T) {
	tests := []struct {
		name        string
		isAmsterdam bool
		gasLimit    uint64
		override    uint64
		isEmpty     bool
		want        uint64
	}{
		{
			name: "pre-amsterdam returns fallback", isAmsterdam: false, gasLimit: 60_000_000,
			override: 0, isEmpty: true, want: BatcherDefaultGasPerTx,
		},
		{
			name: "amsterdam non-empty target", isAmsterdam: true, gasLimit: 60_000_000,
			override: 0, isEmpty: false, want: callRegularGas * 10 / 9, // 13,750
		},
		{
			name: "amsterdam empty 60M (cpsb clamped to floor)", isAmsterdam: true, gasLimit: 60_000_000,
			override: 0, isEmpty: true, want: (callRegularGas + 112*cpsbFloor) * 10 / 9,
		},
		{
			name: "amsterdam empty 100M", isAmsterdam: true, gasLimit: 100_000_000,
			override: 0, isEmpty: true, want: (callRegularGas + 112*1174) * 10 / 9, // 158,336
		},
		{
			name: "override below default is ignored", isAmsterdam: true, gasLimit: 100_000_000,
			override: 30_000, isEmpty: true, want: (callRegularGas + 112*1174) * 10 / 9,
		},
		{
			name: "override above default wins", isAmsterdam: true, gasLimit: 100_000_000,
			override: 200_000, isEmpty: true, want: 200_000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pool := buildTestWalletPool(tt.isAmsterdam, tt.gasLimit, tt.override)
			req := &FundingRequest{IsEmpty: tt.isEmpty}
			require.Equal(t, tt.want, pool.batcherGasFor(req))
		})
	}
}

func TestBatcherBaseGas(t *testing.T) {
	pre := buildTestWalletPool(false, 60_000_000, 0)
	require.Equal(t, uint64(BatcherBaseGas), pre.batcherBaseGas())

	post := buildTestWalletPool(true, 60_000_000, 0)
	require.Equal(t, uint64(65_000), post.batcherBaseGas())
}

func TestPackFundingBatches(t *testing.T) {
	// Helper to build requests of a given emptiness pattern.
	mkReqs := func(flags []bool) []*FundingRequest {
		out := make([]*FundingRequest, len(flags))
		for i, e := range flags {
			out[i] = &FundingRequest{Amount: uint256.NewInt(1), IsEmpty: e}
		}
		return out
	}

	t.Run("empty input produces no batches", func(t *testing.T) {
		pool := buildTestWalletPool(true, 100_000_000, 0)
		batches, total := pool.packFundingBatches(nil)
		require.Empty(t, batches)
		require.Equal(t, uint64(0), total)
	})

	t.Run("small batch fits in one tx", func(t *testing.T) {
		pool := buildTestWalletPool(true, 100_000_000, 0)
		batches, total := pool.packFundingBatches(mkReqs([]bool{true, true, false, true}))
		require.Len(t, batches, 1)
		require.Len(t, batches[0], 4)
		// base + 3 empty + 1 non-empty
		emptyGas := pool.batcherGasFor(&FundingRequest{IsEmpty: true})
		nonEmptyGas := pool.batcherGasFor(&FundingRequest{IsEmpty: false})
		want := pool.batcherBaseGas() + 3*emptyGas + nonEmptyGas
		require.Equal(t, want, total)
	})

	t.Run("batch splits to respect RPC gas cap", func(t *testing.T) {
		pool := buildTestWalletPool(true, 100_000_000, 0)
		// At cpsb=1174, each empty recipient costs ~158k gas. 16M cap lets in ~100.
		// Use 150 empty recipients -> must split across >=2 batches.
		flags := make([]bool, 150)
		for i := range flags {
			flags[i] = true
		}
		batches, _ := pool.packFundingBatches(mkReqs(flags))
		require.GreaterOrEqual(t, len(batches), 2)

		// Each batch (except possibly a singleton first request larger than the cap
		// which doesn't apply here) must fit under BatcherRPCGasCap.
		emptyGas := pool.batcherGasFor(&FundingRequest{IsEmpty: true})
		var seen int
		for _, b := range batches {
			batchGas := pool.batcherBaseGas() + uint64(len(b))*emptyGas
			require.LessOrEqual(t, batchGas, uint64(BatcherRPCGasCap),
				"batch of %d exceeds RPC gas cap", len(b))
			seen += len(b)
		}
		// Recipient count is preserved across batches.
		require.Equal(t, 150, seen)
	})

	t.Run("degrades to legacy batch size pre-amsterdam", func(t *testing.T) {
		pool := buildTestWalletPool(false, 60_000_000, 0)
		// (16M − 50k) / 35k ≈ 455 fits pre-Amsterdam, so 400 fits in one batch.
		flags := make([]bool, 400)
		for i := range flags {
			flags[i] = true
		}
		batches, _ := pool.packFundingBatches(mkReqs(flags))
		require.Len(t, batches, 1)
		require.Len(t, batches[0], 400)
	})
}
