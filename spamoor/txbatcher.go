package spamoor

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethpandaops/spamoor/txbuilder"
	geas "github.com/fjl/geas/asm"
	"github.com/holiman/uint256"
)

// Assembly code for the batcher contract initialization
const batcherGeasInit = `
;; Init code
push @.start
codesize
sub
dup1
push @.start
push0
codecopy
push0
return

.start:
`

// Assembly code for the main batcher contract logic that forwards funds to multiple addresses
const batcherGeasCode = `
push 0                ;; [offset]
jump @loop

loop:
	DUP1              ;; [offset, offset]
	CALLDATASIZE      ;; [calldatasize, offset, offset]
	GT                ;; [calldatasize > offset, offset]
	ISZERO            ;; [calldatasize <= offset, offset]
	JUMPI @exit       ;; [offset]

	;; forward funds
	PUSH 0            ;; [0, offset]
	PUSH 0            ;; [0, 0, offset]
	PUSH 0            ;; [0, 0, 0, offset]
	PUSH 0            ;; [0, 0, 0, 0, offset]

	;; get calldata
	DUP5              ;; [offset, 0, 0, 0, 0, offset]
	CALLDATALOAD      ;; [calldata, 0, 0, 0, 0, offset]

	;; get amount
	PUSH 160          ;; [160, calldata, 0, 0, 0, 0, offset]
	DUP2              ;; [calldata, 160, calldata, 0, 0, 0, 0, offset]
	DUP2              ;; [160, calldata, 160, calldata, 0, 0, 0, 0, offset]
	SHL               ;; [calldata<<160, 160, calldata, 0, 0, 0, 0, offset]
	SWAP1             ;; [160, calldata<<160, calldata, 0, 0, 0, 0, offset]
	SHR               ;; [amount, calldata, 0, 0, 0, 0, offset]
	
	;; get address
	SWAP1             ;; [calldata, amount, 0, 0, 0, 0, offset]
	PUSH 96           ;; [96, calldata, amount, 0, 0, 0, 0, offset]
	SHR               ;; [address, amount, 0, 0, 0, 0, offset]
	
	;; forward funds
	PUSH 30000        ;; [30000, address, amount, 0, 0, 0, 0, offset]
	CALL              ;; [success, offset]
	POP               ;; [offset]

	;; increase offset
	PUSH 32           ;; [32, offset]
	ADD               ;; [offset+32]

	jump @loop

exit:
	SELFBALANCE       ;; [selfbalance]
	DUP1              ;; [selfbalance, selfbalance]
	ISZERO            ;; [selfbalance == 0, selfbalance]
	JUMPI @exit2      ;; [selfbalance]

	;; return leftover funds
	PUSH 0            ;; [0, selfbalance]
	PUSH 0            ;; [0, 0, selfbalance]
	PUSH 0            ;; [0, 0, 0, selfbalance]
	PUSH 0            ;; [0, 0, 0, 0, selfbalance]
	SWAP4             ;; [selfbalance, 0, 0, 0, 0]
	CALLER            ;; [caller, selfbalance, 0, 0, 0, 0]
	GAS               ;; [gas, caller, selfbalance, 0, 0, 0, 0]
	CALL              ;; [success]

exit2:
	STOP
`

const (
	// BatcherBaseGas is the pre-Amsterdam base gas cost for a batcher transaction
	// (tx intrinsic + batcher dispatch overhead). WalletPool.batcherBaseGas uses
	// this value on non-Amsterdam chains and a higher value on Amsterdam to
	// cover the EIP-8037 txpool 110% intrinsic-regular buffer.
	BatcherBaseGas = 50000
	// BatcherDefaultGasPerTx is the pre-Amsterdam default gas cost per recipient
	// in the batch. Used when funding_gas_limit is not configured. On Amsterdam,
	// WalletPool.batcherGasFor computes per-recipient cost dynamically (split by
	// target emptiness + current cpsb) and this constant is the fallback only
	// when txpool.GetCostPerStateByte() returns 0.
	BatcherDefaultGasPerTx = 35000
	// BatcherRPCGasCap is the maximum gas allowed per RPC call (geth default: 16M).
	// Batch boundaries are computed to keep total gas under this cap; see
	// WalletPool.packFundingBatches.
	BatcherRPCGasCap = 16_000_000

	// callRegularGas is a generous estimate of the regular-gas cost a CALL
	// with value incurs inside the batcher contract: CallGasEIP150 (700) +
	// ColdAccountAccess (2600) + CallValueTransfer (9000) + loop opcodes +
	// slack. Used by WalletPool.batcherGasFor.
	callRegularGas = 12_500
)

// TxBatcher manages the deployment and operation of a smart contract that batches
// multiple fund transfers into a single transaction. It compiles and deploys
// assembly code that efficiently forwards funds to multiple recipients.
type TxBatcher struct {
	txpool     *TxPool
	isDeployed bool
	deployMtx  sync.Mutex
	address    common.Address
}

// NewTxBatcher creates a new TxBatcher instance with the specified transaction pool.
// The batcher must be deployed with Deploy() before it can be used.
func NewTxBatcher(txpool *TxPool) *TxBatcher {
	return &TxBatcher{
		txpool: txpool,
	}
}

// Deploy compiles and deploys the batcher smart contract to the blockchain.
// It uses assembly code to create an efficient contract that can forward funds
// to multiple addresses in a single transaction. The deployment is protected
// by a mutex to ensure it only happens once.
//
// Parameters:
//   - ctx: context for the deployment transaction
//   - wallet: wallet to deploy the contract from
//   - client: optional client to use (if nil, uses pool's default client)
func (b *TxBatcher) Deploy(ctx context.Context, wallet *Wallet, client *Client) error {
	b.deployMtx.Lock()
	defer b.deployMtx.Unlock()

	if b.isDeployed {
		return nil
	}

	b.isDeployed = true

	compiler := geas.NewCompiler(nil)

	initcode := compiler.CompileString(batcherGeasInit)
	if initcode == nil {
		return fmt.Errorf("failed to compile initcode")
	}

	batcherGeasCode := compiler.CompileString(batcherGeasCode)
	if batcherGeasCode == nil {
		return fmt.Errorf("failed to compile batcher geas code")
	}

	deployData := append(initcode, batcherGeasCode...)

	if client == nil {
		client = b.txpool.options.ClientPool.GetClient(WithClientSelectionMode(SelectClientByIndex, 0))
		if client == nil {
			return fmt.Errorf("no client available")
		}
	}
	feeCap, tipCap, err := client.GetSuggestedFee(ctx)
	if err != nil {
		return err
	}
	if feeCap.Cmp(big.NewInt(400000000000)) < 0 {
		feeCap = big.NewInt(400000000000)
	}
	if tipCap.Cmp(big.NewInt(200000000000)) < 0 {
		tipCap = big.NewInt(200000000000)
	}

	// Estimate the deploy gas so we stay correct under EIP-8037 where
	// account creation + per-byte code deposit dominate the cost. Fall back
	// to the pre-Amsterdam hard-coded 300k if the RPC path fails.
	deployGas, estErr := client.EstimateGas(ctx, ethereum.CallMsg{
		From:  wallet.GetAddress(),
		To:    nil,
		Value: new(big.Int),
		Data:  deployData,
	})
	if estErr != nil || deployGas == 0 {
		deployGas = fallbackDeployGas(len(deployData), b.txpool.IsAmsterdam(), 0)
	} else if b.txpool.IsAmsterdam() {
		deployGas = deployGas * 12 / 10
	} else {
		deployGas = deployGas * 11 / 10
	}

	txData, err := txbuilder.DynFeeTx(&txbuilder.TxMetadata{
		GasFeeCap: uint256.MustFromBig(feeCap),
		GasTipCap: uint256.MustFromBig(tipCap),
		Gas:       deployGas,
		To:        nil,
		Value:     uint256.NewInt(0),
		Data:      deployData,
	})
	if err != nil {
		return err
	}

	tx, err := wallet.BuildDynamicFeeTx(txData)
	if err != nil {
		return err
	}

	err = b.txpool.SendTransaction(ctx, wallet, tx, &SendTransactionOptions{
		Client:      client,
		Rebroadcast: true,
	})
	if err != nil {
		return err
	}

	b.address = crypto.CreateAddress(wallet.GetAddress(), tx.Nonce())

	return nil
}

// GetAddress returns the deployed contract address.
// Returns zero address if the contract hasn't been deployed yet.
func (b *TxBatcher) GetAddress() common.Address {
	return b.address
}

// GetRequestCalldata encodes funding requests into calldata format expected by the batcher contract.
// Each request is encoded as 32 bytes: 20 bytes for the address and 12 bytes for the amount.
// Returns the encoded calldata that can be used in a transaction to the batcher contract.
func (b *TxBatcher) GetRequestCalldata(requests []*FundingRequest) ([]byte, error) {
	calldata := make([]byte, len(requests)*32)
	for i, request := range requests {
		copy(calldata[i*32:], request.Wallet.GetAddress().Bytes())
		amountBytes := request.Amount.Bytes32()
		copy(calldata[i*32+20:], amountBytes[20:])
	}
	return calldata, nil
}
